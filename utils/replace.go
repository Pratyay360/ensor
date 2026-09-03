package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	FormatBraced = "braced"
	FormatShell  = "shell"
	FormatNode   = "node"
)

func PlaceholderFor(varName, style string) string {
	switch style {
	case FormatShell:
		return "$" + varName
	case FormatNode:
		return "process.env." + varName
	default:
		return "${" + varName + "}"
	}
}

func FindFilesWithSecret(root, secret string) ([]string, error) {
	if secret == "" {
		return nil, nil
	}
	candidates, err := DiscoverFiles(root)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, path := range candidates {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if strings.Contains(string(content), secret) {
			files = append(files, path)
		}
	}
	return files, nil
}

func ReplaceSecretInFiles(root, secret, varName, style, envFile string) (int, int, error) {
	files, err := FindFilesWithSecret(root, secret)
	if err != nil {
		return 0, 0, err
	}

	placeholder := PlaceholderFor(varName, style)
	envAbs, _ := filepath.Abs(envFile)
	var changed, replacements int
	for _, path := range files {
		if abs, err := filepath.Abs(path); err == nil && abs == envAbs {
			continue
		}
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(contentBytes)
		n := strings.Count(content, secret)
		if n == 0 {
			continue
		}
		if err := os.WriteFile(path, []byte(strings.ReplaceAll(content, secret, placeholder)), 0o644); err != nil {
			return changed, replacements, fmt.Errorf("rewriting %s: %w", path, err)
		}
		changed++
		replacements += n
	}
	return changed, replacements, nil
}
