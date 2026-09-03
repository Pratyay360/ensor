package cmd

import (
	"context"
	"math/rand"
	"os"
	"sort"

	log "charm.land/log/v2"
	"github.com/Pratyay360/ensor/utils"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var (
	flagValidate    bool
	flagSkipInvalid bool
	flagEnvFile     string
	flagFormat      string
)

var rootCmd = &cobra.Command{
	Use:   "ensor [directory]",
	Short: "Scan files and folders for secrets and censor them",
	Long:  `Walk the directory recursively, scan files for secrets, censor them + store them in .env`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		RunScan(dir)
	},
}

func Root() *cobra.Command {
	return rootCmd
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func RunScan(dir string) {
	if flagEnvFile == "" {
		flagEnvFile = utils.DotEnvFile
	}
	if _, err := utils.DiscoverFiles(dir); err != nil {
		log.Printf("Error scanning directory: %v\n", err)
		os.Exit(1)
	}
	log.Printf("\n🔍 Scanning directory for secrets using betterleaks...\n")

	findings, err := utils.ScanDirectoryWithBetterleaks(context.Background(), dir, nil, utils.ScanOptions{
		Validate:    flagValidate,
		SkipInvalid: flagSkipInvalid,
	})
	if flagValidate {
		log.Debug("Invalid ones are skipped")
	}

	if err != nil {
		log.Errorf("Error scanning directory with betterleaks: %v", err)
		os.Exit(1)
	}

	log.Infof("Found %d potential secrets", len(findings))
	type secretHit struct {
		file       string
		line       int
		validation string
	}
	var order []string
	hits := map[string]secretHit{}
	for _, finding := range findings {
		if finding.Secret == "" {
			continue
		}
		if _, ok := hits[finding.Secret]; !ok {
			order = append(order, finding.Secret)
			hits[finding.Secret] = secretHit{
				file:       finding.File,
				line:       finding.StartLine,
				validation: finding.ValidationStatus,
			}
		}
	}

	switch flagFormat {
	case utils.FormatBraced, utils.FormatShell, utils.FormatNode:
	default:
		log.Warnf("Unknown --format %q, falling back to %q", flagFormat, utils.FormatBraced)
		flagFormat = utils.FormatBraced
	}

	existingEnv, err := utils.ParseDotEnv(flagEnvFile)
	if err != nil {
		log.Errorf("Error reading %s: %v", flagEnvFile, err)
		os.Exit(1)
	}
	envMap, err := utils.ParseEnvJSON(utils.EnvJSONFile)
	if err != nil {
		log.Errorf("Error reading %s: %v", utils.EnvJSONFile, err)
		os.Exit(1)
	}
	secretToName := map[string]string{}
	for name, val := range envMap {
		if _, ok := secretToName[val]; !ok {
			secretToName[val] = name
		}
	}
	for name, val := range existingEnv {
		if _, ok := secretToName[val]; !ok {
			secretToName[val] = name
		}
	}

	pending := map[string]string{}
	taken := func(name string) bool {
		if _, ok := envMap[name]; ok {
			return true
		}
		if _, ok := existingEnv[name]; ok {
			return true
		}
		if _, ok := pending[name]; ok {
			return true
		}
		return false
	}

	type redaction struct{ name, secret string }
	var redactions []redaction

	randomName := func() string {
		letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := make([]rune, 6)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return "SECRET_" + string(b)
	}

	for _, secret := range order {
		hit := hits[secret]

		if name, ok := secretToName[secret]; ok {
			log.Infof("✅ Known secret in %s (line %d) maps to %s", hit.file, hit.line, utils.PlaceholderFor(name, flagFormat))
			redactions = append(redactions, redaction{name, secret})
			continue
		}

		redact, errPrompt := utils.PromptRedactOrSkip(hit.file, hit.line, secret, hit.validation)
		if errPrompt != nil || !redact {
			log.Infof("⏭️  Skipped secret in %s (line %d)", hit.file, hit.line)
			continue
		}
		suggestion := randomName()
		name, err := utils.PromptForSecretNameWithDefault(suggestion, taken)
		if err != nil {
			log.Infof("⏭️  Skipped secret in %s (line %d) (prompt cancelled)", hit.file, hit.line)
			continue
		}

		pending[name] = secret
		secretToName[secret] = name
		redactions = append(redactions, redaction{name, secret})
		log.Infof("📝 Queued %s for %s", name, hit.file)
	}
	if len(pending) > 0 {
		if err := utils.AppendToDotEnv(flagEnvFile, pending); err != nil {
			log.Errorf("❌ Error saving to %s: %v", flagEnvFile, err)
			os.Exit(1)
		}
		log.Infof("💾 Saved %d variable(s) to %s", len(pending), flagEnvFile)
	}

	sort.Slice(redactions, func(i, j int) bool { return len(redactions[i].secret) > len(redactions[j].secret) })
	for _, r := range redactions {
		changed, n, err := utils.ReplaceSecretInFiles(dir, r.secret, r.name, flagFormat, flagEnvFile)
		if err != nil {
			log.Errorf("❌ Error replacing secret for %s: %v", r.name, err)
			continue
		}
		if changed > 0 {
			log.Infof("✅ Replaced %d occurrence(s) with %s in %d file(s)", n, utils.PlaceholderFor(r.name, flagFormat), changed)
		}
	}

	log.Info("Scan and redaction complete!")
}

var redactCmd = &cobra.Command{
	Use:   "",
	Short: "Scan files for secrets",
	Long: `Scan all files recursively in the specified directory (defaulting to ".") for secrets.
Prompts the user to give a variable name for each detected secret, saves the secret key-value
pair to the env.json file, and replaces the secret in-place inside the source file with a placeholder.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		RunScan(dir)
	},
}

func init() {
	rootCmd.AddCommand(redactCmd)
	rootCmd.Flags().BoolVar(&flagValidate, "validate", false, "verify secrets against their providers to cut false positives (requires network access)")
	rootCmd.Flags().BoolVar(&flagSkipInvalid, "skip-invalid", true, "with --validate, ignore secrets the provider confirmed as invalid")
	rootCmd.Flags().StringVar(&flagEnvFile, "env-file", utils.DotEnvFile, "dotenv file to append new variables to (never overwrites)")
	rootCmd.Flags().StringVar(&flagFormat, "format", utils.FormatBraced, "replacement placeholder style: braced (${NAME}), shell ($NAME) or node (process.env.NAME)")
	redactCmd.Flags().BoolVar(&flagValidate, "validate", false, "verify secrets against their providers to cut false positives (requires network access)")
	redactCmd.Flags().BoolVar(&flagSkipInvalid, "skip-invalid", true, "with --validate, ignore secrets the provider confirmed as invalid")
	redactCmd.Flags().StringVar(&flagEnvFile, "env-file", utils.DotEnvFile, "dotenv file to append new variables to (never overwrites)")
	redactCmd.Flags().StringVar(&flagFormat, "format", utils.FormatBraced, "replacement placeholder style: braced (${NAME}), shell ($NAME) or node (process.env.NAME)")
}
