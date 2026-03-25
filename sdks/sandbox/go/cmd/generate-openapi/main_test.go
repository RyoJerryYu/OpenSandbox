package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratorSpecsExist(t *testing.T) {
	repoRoot, err := repoRootFromWorkingDir()
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	if err := validateSpecsExist(repoRoot); err != nil {
		t.Fatalf("expected required specs to exist: %v", err)
	}
}

func TestGenerationTargetsHaveExpectedSpecsAndOutputs(t *testing.T) {
	repoRoot, err := repoRootFromWorkingDir()
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	targets := generationTargets(repoRoot)
	if len(targets) != 3 {
		t.Fatalf("expected 3 generation targets, got %d", len(targets))
	}

	expected := map[string]string{
		"lifecycle": filepath.Join(repoRoot, "specs", "sandbox-lifecycle.yml"),
		"execd":     filepath.Join(repoRoot, "specs", "execd-api.yaml"),
		"egress":    filepath.Join(repoRoot, "specs", "egress-api.yaml"),
	}

	for _, target := range targets {
		wantSpec, ok := expected[target.name]
		if !ok {
			t.Fatalf("unexpected target name: %s", target.name)
		}
		if target.specPath != wantSpec {
			t.Fatalf("target %s spec mismatch: got %s want %s", target.name, target.specPath, wantSpec)
		}
		wantOut := filepath.Join(repoRoot, "sdks", "sandbox", "go", "sandbox", "internal", "openapi", target.name)
		if target.outputDir != wantOut {
			t.Fatalf("target %s output mismatch: got %s want %s", target.name, target.outputDir, wantOut)
		}
		if target.packageName == "" {
			t.Fatalf("target %s package name should not be empty", target.name)
		}
	}
}

func TestRewriteOpenAPI31NullablesToNullable(t *testing.T) {
	input := strings.Join([]string{
		"openapi: 3.1.0",
		"components:",
		"  schemas:",
		"    Example:",
		"      properties:",
		"        timeout:",
		"          oneOf:",
		"            - type: integer",
		"              minimum: 60",
		"            - type: 'null'",
		"        expiresAt:",
		"          oneOf:",
		"            - type: string",
		"              format: date-time",
		"            - type: 'null'",
	}, "\n")

	output := rewriteOpenAPI31NullablesToNullable(input)

	if strings.Contains(output, "openapi: 3.1.0") {
		t.Fatal("expected OpenAPI version to be downgraded for generation")
	}
	if !strings.Contains(output, "openapi: 3.0.3") {
		t.Fatal("expected OpenAPI 3.0.3 in rewritten spec")
	}
	if strings.Contains(output, "- type: 'null'") {
		t.Fatal("expected null branch to be removed")
	}
	if !strings.Contains(output, "nullable: true") {
		t.Fatal("expected nullable flag to be introduced")
	}
}

func TestPrepareSpecForGenerationRewritesLifecycleSpecOnly(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "sandbox-lifecycle.yml")
	specContent := strings.Join([]string{
		"openapi: 3.1.0",
		"components:",
		"  schemas:",
		"    Example:",
		"      properties:",
		"        timeout:",
		"          oneOf:",
		"            - type: integer",
		"            - type: 'null'",
	}, "\n")
	if err := os.WriteFile(specPath, []byte(specContent), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	target := generationTarget{
		name:      "lifecycle",
		specPath:  specPath,
		outputDir: filepath.Join(dir, "out"),
	}
	preparedPath, cleanup, err := prepareSpecForGeneration(target)
	if err != nil {
		t.Fatalf("prepare spec: %v", err)
	}
	defer cleanup()

	if preparedPath == specPath {
		t.Fatal("expected lifecycle spec to be rewritten to a temp file")
	}

	rewritten, err := os.ReadFile(preparedPath)
	if err != nil {
		t.Fatalf("read rewritten spec: %v", err)
	}
	if !strings.Contains(string(rewritten), "nullable: true") {
		t.Fatal("expected rewritten lifecycle spec to contain nullable flag")
	}

	otherTarget := generationTarget{name: "execd", specPath: specPath}
	otherPath, otherCleanup, err := prepareSpecForGeneration(otherTarget)
	if err != nil {
		t.Fatalf("prepare non-lifecycle spec: %v", err)
	}
	defer otherCleanup()
	if otherPath != specPath {
		t.Fatal("expected non-lifecycle spec to be used as-is")
	}
}
