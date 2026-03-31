package services

import (
	"context"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

// SandboxFiles defines high-level filesystem operations against a connected sandbox.
type SandboxFiles interface {
	GetFileInfo(ctx context.Context, paths []string) (map[string]models.EntryInfo, error)
	Search(ctx context.Context, entry models.SearchEntry) ([]models.EntryInfo, error)
	CreateDirectories(ctx context.Context, entries []models.WriteEntry) error
	DeleteDirectories(ctx context.Context, paths []string) error
	WriteFiles(ctx context.Context, entries []models.WriteEntry) error
	ReadFile(ctx context.Context, path string, opts *models.ReadFileOptions) (string, error)
	ReadBytes(ctx context.Context, path string, opts *models.ReadFileOptions) ([]byte, error)
	DeleteFiles(ctx context.Context, paths []string) error
	MoveFiles(ctx context.Context, entries []models.MoveEntry) error
	ReplaceContents(ctx context.Context, entries []models.ContentReplaceEntry) error
	SetPermissions(ctx context.Context, entries []models.SetPermissionEntry) error
}
