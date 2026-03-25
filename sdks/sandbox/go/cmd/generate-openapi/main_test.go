package main

import "testing"

func TestGeneratorSpecsExist(t *testing.T) {
	repoRoot, err := repoRootFromWorkingDir()
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	if err := validateSpecsExist(repoRoot); err != nil {
		t.Fatalf("expected required specs to exist: %v", err)
	}
}
