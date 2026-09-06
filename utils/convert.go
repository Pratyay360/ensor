package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"go.yaml.in/yaml/v3"
)

// EnvFormat is a flat key=value file format supported by convert.
type EnvFormat string

const (
	FormatEnv  EnvFormat = "env"
	FormatJSON EnvFormat = "json"
	FormatYAML EnvFormat = "yaml"
	FormatTOML EnvFormat = "toml"
)

// NormalizeEnvFormat maps a --from/--to flag value to an EnvFormat.
func NormalizeEnvFormat(s string) (EnvFormat, error) {
	switch EnvFormat(strings.ToLower(strings.TrimSpace(s))) {
	case FormatEnv, "dotenv":
		return FormatEnv, nil
	case FormatJSON:
		return FormatJSON, nil
	case FormatYAML, "yml":
		return FormatYAML, nil
	case FormatTOML:
		return FormatTOML, nil
	default:
		return "", fmt.Errorf("unknown format %q: want env, json, yaml or toml", s)
	}
}

// DetectEnvFormat infers the format from the file extension.
func DetectEnvFormat(path string) (EnvFormat, error) {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".json"):
		return FormatJSON, nil
	case strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml"):
		return FormatYAML, nil
	case strings.HasSuffix(lower, ".toml"):
		return FormatTOML, nil
	case strings.HasSuffix(lower, ".env"):
		return FormatEnv, nil
	default:
		return "", fmt.Errorf("cannot infer format from %q: pass --from env|json|yaml|toml", path)
	}
}

// ParseEnvFile reads path as a flat string map in the given format.
func ParseEnvFile(path string, format EnvFormat) (map[string]string, error) {
	switch format {
	case FormatEnv:
		return ParseDotEnv(path)
	case FormatJSON:
		return parseFlatUnmarshal(path, json.Unmarshal)
	case FormatYAML:
		return parseFlatUnmarshal(path, yaml.Unmarshal)
	case FormatTOML:
		return parseFlatUnmarshal(path, toml.Unmarshal)
	default:
		return nil, fmt.Errorf("unknown format %q: want env, json, yaml or toml", format)
	}
}

func parseFlatUnmarshal(path string, unmarshal func([]byte, any) error) (map[string]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(b))) == 0 {
		return map[string]string{}, nil
	}
	raw := map[string]any{}
	if err := unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		if v == nil {
			out[k] = ""
			continue
		}
		if s, ok := v.(string); ok {
			out[k] = s
			continue
		}
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}

// MarshalEnv encodes a flat string map in the given format.
func MarshalEnv(data map[string]string, format EnvFormat) ([]byte, error) {
	switch format {
	case FormatEnv:
		return marshalDotEnv(data), nil
	case FormatJSON:
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, err
		}
		return append(b, '\n'), nil
	case FormatYAML:
		return yaml.Marshal(data)
	case FormatTOML:
		return toml.Marshal(data)
	default:
		return nil, fmt.Errorf("unknown format %q: want env, json, yaml or toml", format)
	}
}

func marshalDotEnv(data map[string]string) []byte {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		v := data[k]
		if strings.ContainsAny(v, " \t#\"'") || strings.Contains(v, "\n") {
			escaped := strings.ReplaceAll(v, `\`, `\\`)
			escaped = strings.ReplaceAll(escaped, `"`, `\"`)
			v = `"` + escaped + `"`
		}
		fmt.Fprintf(&b, "%s=%s\n", k, v)
	}
	return []byte(b.String())
}
