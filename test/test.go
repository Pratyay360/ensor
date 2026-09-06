// A simple, testcase generator ..
// glorified uuid generator, really

package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

type Secret struct {
	Name  string
	Value string
}

func main() {
	const (
		secretCount = 10
		root        = "testdata"
	)
	if err := os.RemoveAll(root); err != nil {
		panic(err)
	}

	must(os.MkdirAll(root, 0755))
	must(os.MkdirAll(filepath.Join(root, "src"), 0755))
	must(os.MkdirAll(filepath.Join(root, "config"), 0755))
	must(os.MkdirAll(filepath.Join(root, "scripts"), 0755))

	secrets := make([]Secret, 0, secretCount)

	for i := 0; i < secretCount; i++ {
		secrets = append(secrets, Secret{
			Name:  fmt.Sprintf("SECRET_%04d", i),
			Value: randomSecret(32),
		})
	}

	// Generate files.
	generateEnv(root, secrets)
	generateJSON(root, secrets)
	generateYAML(root, secrets)
	generateGo(root, secrets)
	generateShell(root, secrets)

	// Generate duplicates in another file.
	generateDuplicates(root, secrets)

	fmt.Printf(
		"Generated %d fake secrets in %s\n",
		len(secrets),
		root,
	)
}

func generateEnv(root string, secrets []Secret) {
	var b strings.Builder

	for _, s := range secrets {
		fmt.Fprintf(
			&b,
			"%s=%s\n",
			s.Name,
			s.Value,
		)
	}

	must(os.WriteFile(
		filepath.Join(root, ".env"),
		[]byte(b.String()),
		0644,
	))
}

func generateJSON(root string, secrets []Secret) {
	var b strings.Builder

	b.WriteString("{\n")

	for i, s := range secrets {
		fmt.Fprintf(
			&b,
			"  %q: %q",
			s.Name,
			s.Value,
		)

		if i < len(secrets)-1 {
			b.WriteString(",")
		}

		b.WriteString("\n")
	}

	b.WriteString("}\n")

	must(os.WriteFile(
		filepath.Join(root, "config", "secrets.json"),
		[]byte(b.String()),
		0644,
	))
}

func generateYAML(root string, secrets []Secret) {
	var b strings.Builder

	for _, s := range secrets {
		fmt.Fprintf(
			&b,
			"%s: %s\n",
			s.Name,
			s.Value,
		)
	}

	must(os.WriteFile(
		filepath.Join(root, "config", "secrets.yaml"),
		[]byte(b.String()),
		0644,
	))
}

func generateGo(root string, secrets []Secret) {
	var b strings.Builder

	b.WriteString("package config\n\nconst (\n")

	for _, s := range secrets {
		fmt.Fprintf(
			&b,
			"\t%s = %q\n",
			s.Name,
			s.Value,
		)
	}

	b.WriteString(")\n")

	must(os.WriteFile(
		filepath.Join(root, "src", "config.go"),
		[]byte(b.String()),
		0644,
	))
}

func generateShell(root string, secrets []Secret) {
	var b strings.Builder

	b.WriteString("#!/bin/sh\n\n")

	for _, s := range secrets {
		fmt.Fprintf(
			&b,
			"export %s='%s'\n",
			s.Name,
			s.Value,
		)
	}

	must(os.WriteFile(
		filepath.Join(root, "scripts", "config.sh"),
		[]byte(b.String()),
		0755,
	))
}

func generateDuplicates(root string, secrets []Secret) {
	var b strings.Builder

	for i, s := range secrets {
		if i%10 != 0 {
			continue
		}

		fmt.Fprintf(
			&b,
			"PRIMARY_%04d=%s\n",
			i,
			s.Value,
		)

		fmt.Fprintf(
			&b,
			"BACKUP_%04d=%s\n",
			i,
			s.Value,
		)

		fmt.Fprintf(
			&b,
			"LEGACY_%04d=%s\n",
			i,
			s.Value,
		)

		fmt.Fprintf(
			&b,
			"THIRD_PARTY_%04d=%s\n",
			i,
			s.Value,
		)

		b.WriteString("\n")
	}

	must(os.WriteFile(
		filepath.Join(root, "duplicates.env"),
		[]byte(b.String()),
		0644,
	))
}

func randomSecret(n int) string {
	buf := make([]byte, n)

	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	return hex.EncodeToString(buf)
}

func randomInt(max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}

	return n.Int64()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
