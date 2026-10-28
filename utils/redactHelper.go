package utils

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	lipgloss "charm.land/lipgloss/v2"
	log "charm.land/log/v2"
	"github.com/charmbracelet/x/term"
)

func IsInteractive() bool {
	return term.IsTerminal(os.Stdin.Fd())
}

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

const (
	ActionRedact = "redact"
	ActionSkip   = "skip"
	ActionQuit   = "quit"
)

type hotkeySelect struct {
	*huh.Select[string]
	choice *string
}

func newHotkeySelect(choice *string) *hotkeySelect {
	return &hotkeySelect{
		Select: huh.NewSelect[string](),
		choice: choice,
	}
}

func (h *hotkeySelect) Update(msg tea.Msg) (huh.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch strings.ToLower(keyMsg.String()) {
		case "r":
			*h.choice = ActionRedact
			return h, huh.NextField
		case "s":
			*h.choice = ActionSkip
			return h, huh.NextField
		case "q":
			*h.choice = ActionQuit
			return h, huh.NextField
		}
	}
	m, cmd := h.Select.Update(msg)
	if sel, ok := m.(*huh.Select[string]); ok {
		h.Select = sel
	}
	return h, cmd
}

func (h *hotkeySelect) GetValue() any {
	if h.choice != nil && *h.choice != "" {
		return *h.choice
	}
	return h.Select.GetValue()
}

func (h *hotkeySelect) WithTheme(t huh.Theme) huh.Field {
	h.Select.WithTheme(t)
	return h
}

func (h *hotkeySelect) WithKeyMap(k *huh.KeyMap) huh.Field {
	h.Select.WithKeyMap(k)
	return h
}

func (h *hotkeySelect) WithWidth(w int) huh.Field {
	h.Select.WithWidth(w)
	return h
}

func (h *hotkeySelect) WithHeight(ht int) huh.Field {
	h.Select.WithHeight(ht)
	return h
}

func (h *hotkeySelect) WithPosition(p huh.FieldPosition) huh.Field {
	h.Select.WithPosition(p)
	return h
}

func PromptRedactOrSkip(filename string, line int, secretValue string, validationStatus string) (string, error) {
	renderSecretBox(filename, line, secretValue, validationStatus)

	var choice string
	sel := newHotkeySelect(&choice)
	sel.
		Title("Redact or skip?").
		Description("Choose whether to redact this secret, skip it, or quit").
		Options(
			huh.NewOption("Redact (r)", ActionRedact),
			huh.NewOption("Skip (s)", ActionSkip),
			huh.NewOption("Quit (q)", ActionQuit),
		).
		Value(&choice)

	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))

	form := huh.NewForm(
		huh.NewGroup(sel),
	).WithTheme(customTheme{}).WithKeyMap(km)

	if err := form.Run(); err != nil {
		return "", err
	}
	return choice, nil
}

func PromptForSecretNameWithDefault(suggestion string, exists func(string) bool) (string, error) {
	for {
		var input string
		km := huh.NewDefaultKeyMap()
		km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))
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
		).WithTheme(customTheme{}).WithKeyMap(km)

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
