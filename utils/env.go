package utils

import (
	"bufio"
	"encoding/json"
	log "charm.land/log/v2"
	"fmt"
	"os"
	"sort"
	"strings"
)

const DotEnvFile = ".env"

const EnvJSONFile = "env.json"

// unescapeDoubleQuoted resolves the \\ and \" escapes inside a
// double-quoted dotenv value. Any other backslash sequence is kept as-is.
func unescapeDoubleQuoted(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	escaped := false
	for _, r := range s {
		if escaped {
			if r != '"' && r != '\\' {
				b.WriteByte('\\')
			}
			b.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		b.WriteRune(r)
	}
	if escaped {
		b.WriteByte('\\')
	}
	return b.String()
}

func ParseDotEnv(path string) (map[string]string, error) {
	if path == "" {
		path = DotEnvFile
	}
	m := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return nil, err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx >= 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			// Strip optional surrounding quotes.
			if len(val) >= 2 {
				if val[0] == '"' && val[len(val)-1] == '"' {
					val = unescapeDoubleQuoted(val[1 : len(val)-1])
				} else if val[0] == '\'' && val[len(val)-1] == '\'' {
					val = val[1 : len(val)-1]
				}
			}
			if key != "" {
				m[key] = val
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func ParseEnvJSON(path string) (map[string]string, error) {
	if path == "" {
		path = EnvJSONFile
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	// Empty file
	if len(strings.TrimSpace(string(b))) == 0 {
		return map[string]string{}, nil
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		log.Error("parsing env json", "path", path, "err", err)
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if m == nil {
		m = map[string]string{}
	}
	return m, nil
}

func AppendToDotEnv(path string, pending map[string]string) error {
	if len(pending) == 0 {
		return nil
	}
	if path == "" {
		path = DotEnvFile
	}

	existing, err := ParseDotEnv(path)
	if err != nil {
		return err
	}

	filtered := map[string]string{}
	for k, v := range pending {
		if _, ok := existing[k]; !ok {
			filtered[k] = v
		}
	}
	if len(filtered) == 0 {
		return nil
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// If file exists and is non-empty and doesn't end with newline, add one.
	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		b, err := os.ReadFile(path)
		if err == nil && len(b) > 0 && b[len(b)-1] != '\n' {
			if _, err := f.WriteString("\n"); err != nil {
				return err
			}
		}
	}

	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bufio.NewWriter(f)
	for _, k := range keys {
		v := filtered[k]
		// Quote value if it contains spaces, #, or quotes/backslashes.
		needsQuote := strings.ContainsAny(v, " \t#\"'") || strings.Contains(v, "\n")
		if needsQuote {
			escaped := strings.ReplaceAll(v, `\`, `\\`)
			escaped = strings.ReplaceAll(escaped, `"`, `\"`)
			v = `"` + escaped + `"`
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return w.Flush()
}
