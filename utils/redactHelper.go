package utils

import (
	"fmt"
	"os"
	"strings"

	huh "charm.land/huh/v2"
	lipgloss "charm.land/lipgloss/v2"
	log "charm.land/log/v2"
)

func SecretOnLine(filePath string, lineNum int, secretValue string) (bool, error) {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	lines := strings.Split(string(contentBytes), "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return false, nil
	}
	return strings.Contains(lines[lineNum-1], secretValue), nil
}

func RedactSecretOnLine(filePath string, lineNum int, secretValue string, varName string) error {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(contentBytes)
	lines := strings.Split(content, "\n")

	if lineNum < 1 || lineNum > len(lines) {
		log.Error("line number out of range", "line", lineNum, "file", filePath)
	}

	replacement := "${" + varName + "}"
	lines[lineNum-1] = strings.ReplaceAll(lines[lineNum-1], secretValue, replacement)

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0o644)
}

func PromptForSecretValue(varName string) (string, error) {
	var value string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Enter value for ${%s}:", varName)).
				Value(&value),
		),
	).WithTheme(huh.ThemeFunc(huh.ThemeCharm))

	if err := form.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(value), nil
}

type customTheme struct{}

func (c customTheme) Theme(isDark bool) *huh.Styles {
	s := huh.ThemeCharm(isDark)
	s.Focused.FocusedButton = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#33d17a")).
		Bold(true).
		Padding(0, 3).
		MarginRight(2)
	s.Focused.BlurredButton = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 3).
		MarginRight(2)
	return s
}

func renderSecretBox(filename string, line int, secretValue string, validationStatus string) {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#33d17a"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Width(12)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
	secretStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#dc8add")).Bold(true)
	validationStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA657")).Italic(true)

	header := headerStyle.Render("🔍 Secrets found")
	fileRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("File:"), valueStyle.Render(filename))
	lineRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Line:"), valueStyle.Render(fmt.Sprintf("%d", line)))

	boxLines := []string{header, fileRow, lineRow}
	if validationStatus != "" {
		validationRow := lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Validation:"), validationStyle.Render(validationStatus))
		boxLines = append(boxLines, validationRow)
	}
	boxLines = append(boxLines, "", secretStyle.Render(secretValue))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffa657")).
		Padding(1, 2).
		Width(60).
		Render(lipgloss.JoinVertical(lipgloss.Left, boxLines...))
	fmt.Println(box)
}

func PromptRedactOrSkip(filename string, line int, secretValue string, validationStatus string) (bool, error) {
	renderSecretBox(filename, line, secretValue, validationStatus)

	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Redact or skip?").
				Description("Choose whether to redact this secret or skip it").
				Options(
					huh.NewOption("Redact (r)", "redact"),
					huh.NewOption("Skip (s)", "skip"),
				).
				Value(&choice),
		),
	).WithTheme(customTheme{})

	if err := form.Run(); err != nil {
		return false, err
	}
	return choice == "redact", nil
}

// PromptForSecretNameWithDefault prompts "Enter SECRET Name (default:
// <suggestion>)" and returns the chosen name. Pressing Enter accepts the
// suggestion. Empty, malformed, or already-taken names are rejected with a
// reprompt; exists reports whether a name is already used.
func PromptForSecretNameWithDefault(suggestion string, exists func(string) bool) (string, error) {
	for {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(fmt.Sprintf("Enter SECRET Name (default: %s):", suggestion)).
					Placeholder(suggestion).
					Value(&input).
					Validate(func(s string) error {
						name := strings.TrimSpace(s)
						if name == "" {
							name = suggestion
						}
						if exists(name) {
							return fmt.Errorf("%q is already in use", name)
						}
						return nil
					}),
			),
		).WithTheme(customTheme{})

		if err := form.Run(); err != nil {
			return "", err
		}
		name := strings.TrimSpace(input)
		if name == "" {
			name = suggestion
		}
		if exists(name) {
			log.Warnf("⚠️  %s is already in use, choose another name.", name)
			continue
		}
		return name, nil
	}
}

func PromptConfirmSave(filename string, line int, secretValue string, validationStatus string) (bool, error) {
	renderSecretBox(filename, line, secretValue, validationStatus)

	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to redact and save this secret?").
				Affirmative("Redact").
				Negative("Skip").
				Value(&confirm),
		),
	).WithTheme(customTheme{})

	if err := form.Run(); err != nil {
		return false, err
	}
	return confirm, nil
}
