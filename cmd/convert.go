package cmd

import (
	"fmt"
	"os"

	log "charm.land/log/v2"
	"github.com/Pratyay360/ensor/utils"
	"github.com/spf13/cobra"
)

var (
	convertFrom string
	convertTo   string
)

// convertCmd converts a flat key=value file between dotenv, JSON, YAML
// and TOML, in any direction.
var convertCmd = &cobra.Command{
	Use:   "convert [input] [output]",
	Short: "Convert between .env, JSON, YAML and TOML",
	Long: `Convert a flat key=value file between dotenv, JSON, YAML and TOML,
in any direction. Formats are inferred from the file extensions unless
--from/--to is given. With no output file the result goes to stdout.
Examples:
  ensor convert .env env.json
  ensor convert env.json env.yaml
  ensor convert config.yaml .env
  ensor convert --to toml .env > env.toml`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := utils.DotEnvFile
		if len(args) > 0 {
			input = args[0]
		}

		from, err := resolveConvertFormat(convertFrom, input, "--from")
		if err != nil {
			return err
		}
		data, err := utils.ParseEnvFile(input, from)
		if err != nil {
			return err
		}

		var to utils.EnvFormat
		if len(args) > 1 {
			to, err = resolveConvertFormat(convertTo, args[1], "--to")
			if err != nil {
				return err
			}
		} else if convertTo != "" {
			to, err = utils.NormalizeEnvFormat(convertTo)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("no output file and no --to flag: pass an output file or --to env|json|yaml|toml")
		}

		out, err := utils.MarshalEnv(data, to)
		if err != nil {
			return err
		}

		if len(args) < 2 {
			_, err = fmt.Fprint(cmd.OutOrStdout(), string(out))
			return err
		}
		if err := os.WriteFile(args[1], out, 0o644); err != nil {
			return err
		}
		log.Infof("Converted %s (%s) to %s (%s)", input, from, args[1], to)
		return nil
	},
}

func resolveConvertFormat(flag, path, name string) (utils.EnvFormat, error) {
	if flag != "" {
		return utils.NormalizeEnvFormat(flag)
	}
	return utils.DetectEnvFormat(path)
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVar(&convertFrom, "from", "", "input format: env, json, yaml or toml (default: inferred from input extension)")
	convertCmd.Flags().StringVar(&convertTo, "to", "", "output format: env, json, yaml or toml (default: inferred from output extension)")
}
