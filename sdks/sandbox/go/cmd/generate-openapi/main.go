package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func rewriteOpenAPI31NullablesToNullable(input string) string {
	output := strings.Replace(input, "openapi: 3.1.0", "openapi: 3.0.3", 1)

	replacements := []struct {
		old string
		new string
	}{
		{
			old: strings.Join([]string{
				"          oneOf:",
				"            - type: string",
				"              format: date-time",
				"            - type: 'null'",
			}, "\n"),
			new: strings.Join([]string{
				"          type: string",
				"          format: date-time",
				"          nullable: true",
			}, "\n"),
		},
		{
			old: strings.Join([]string{
				"          oneOf:",
				"            - type: integer",
				"              minimum: 60",
				"            - type: 'null'",
			}, "\n"),
			new: strings.Join([]string{
				"          type: integer",
				"          minimum: 60",
				"          nullable: true",
			}, "\n"),
		},
		{
			old: strings.Join([]string{
				"          oneOf:",
				"            - type: integer",
				"            - type: 'null'",
			}, "\n"),
			new: strings.Join([]string{
				"          type: integer",
				"          nullable: true",
			}, "\n"),
		},
	}

	for _, replacement := range replacements {
		output = strings.ReplaceAll(output, replacement.old, replacement.new)
	}

	return output
}

func prepareSpecForGeneration(target generationTarget) (string, func(), error) {
	if target.name != "lifecycle" {
		return target.specPath, func() {}, nil
	}

	content, err := os.ReadFile(target.specPath)
	if err != nil {
		return "", nil, fmt.Errorf("read %s spec: %w", target.name, err)
	}

	rewritten := rewriteOpenAPI31NullablesToNullable(string(content))
	tmpFile, err := os.CreateTemp("", target.name+"-*.yaml")
	if err != nil {
		return "", nil, fmt.Errorf("create temp spec for %s: %w", target.name, err)
	}
	if _, err := tmpFile.WriteString(rewritten); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("write temp spec for %s: %w", target.name, err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("close temp spec for %s: %w", target.name, err)
	}

	cleanup := func() {
		_ = os.Remove(tmpFile.Name())
	}
	return tmpFile.Name(), cleanup, nil
}

func runGeneration(target generationTarget) error {
	if err := os.MkdirAll(target.outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir for %s: %w", target.name, err)
	}

	specPath, cleanup, err := prepareSpecForGeneration(target)
	if err != nil {
		return err
	}
	defer cleanup()

	outputFile := filepath.Join(target.outputDir, "client.gen.go")
	cmd := exec.Command(
		"oapi-codegen",
		"-generate", "types,client",
		"-package", target.packageName,
		"-o", outputFile,
		specPath,
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
