package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func repoRootFromWorkingDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		specsDir := filepath.Join(dir, "specs")
		info, err := os.Stat(specsDir)
		if err == nil && info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repo root from %s", dir)
		}
		dir = parent
	}
}

func requiredSpecs(repoRoot string) []string {
	return []string{
		filepath.Join(repoRoot, "specs", "sandbox-lifecycle.yml"),
		filepath.Join(repoRoot, "specs", "execd-api.yaml"),
		filepath.Join(repoRoot, "specs", "egress-api.yaml"),
	}
}

func validateSpecsExist(repoRoot string) error {
	for _, spec := range requiredSpecs(repoRoot) {
		if _, err := os.Stat(spec); err != nil {
			return fmt.Errorf("missing spec %s: %w", spec, err)
		}
	}
	return nil
}

func main() {
	repoRoot, err := repoRootFromWorkingDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve repo root: %v\n", err)
		os.Exit(1)
	}

	if err := validateSpecsExist(repoRoot); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("OpenAPI generator bootstrap is ready.")
}
