package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func IgnoreFile(path string, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return false
	}

	parts := strings.Split(rel, string(filepath.Separator))
	base := strings.ToLower(parts[len(parts)-1])
	if base == ".env" || base == "env.json" {
		return true
	}

	for _, part := range parts {
		if strings.HasPrefix(strings.ToLower(part), ".") {
			return true
		}
	}

	return false
}

func DiscoverFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Not being able to inspect the root itself is fatal, but any
			// other unreadable entry (permissions, races, ...) should be
			// skipped so a single entry can't abort the whole walk.
			if path == root {
				return err
			}
			if info != nil && info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if IgnoreFile(path, root) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
