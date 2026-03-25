package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type generationTarget struct {
	name        string
	specPath    string
	outputDir   string
	packageName string
}

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
	targets := generationTargets(repoRoot)
	specs := make([]string, 0, len(targets))
	for _, target := range targets {
		specs = append(specs, target.specPath)
	}
	return specs
}

func generationTargets(repoRoot string) []generationTarget {
	openapiRoot := filepath.Join(repoRoot, "sdks", "sandbox", "go", "sandbox", "internal", "openapi")
	return []generationTarget{
		{
			name:        "lifecycle",
			specPath:    filepath.Join(repoRoot, "specs", "sandbox-lifecycle.yml"),
			outputDir:   filepath.Join(openapiRoot, "lifecycle"),
			packageName: "lifecycle",
		},
		{
			name:        "execd",
			specPath:    filepath.Join(repoRoot, "specs", "execd-api.yaml"),
			outputDir:   filepath.Join(openapiRoot, "execd"),
			packageName: "execd",
		},
		{
			name:        "egress",
			specPath:    filepath.Join(repoRoot, "specs", "egress-api.yaml"),
			outputDir:   filepath.Join(openapiRoot, "egress"),
			packageName: "egress",
		},
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

func runGeneration(target generationTarget) error {
	if err := os.MkdirAll(target.outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir for %s: %w", target.name, err)
	}

	outputFile := filepath.Join(target.outputDir, "client.gen.go")
	cmd := exec.Command(
		"oapi-codegen",
		"-generate", "types,client",
		"-package", target.packageName,
		"-o", outputFile,
		target.specPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generate %s client: %w", target.name, err)
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

	for _, target := range generationTargets(repoRoot) {
		if err := runGeneration(target); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Println("OpenAPI client generation completed.")
}
