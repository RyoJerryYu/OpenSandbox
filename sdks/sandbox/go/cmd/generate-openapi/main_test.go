package main

import (
	"path/filepath"
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
