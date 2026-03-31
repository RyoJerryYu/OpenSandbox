package models

import "time"

// EntryInfo contains filesystem metadata for one sandbox path.
type EntryInfo struct {
	Path       string
	Size       int64
	ModifiedAt time.Time
	CreatedAt  time.Time
	Mode       int
	Owner      string
	Group      string
}

// WriteEntry describes a file or directory write operation.
type WriteEntry struct {
	Path     string
	Data     any
	Mode     int
	Owner    string
	Group    string
	Encoding string
}

// SearchEntry describes a filesystem search request.
type SearchEntry struct {
	Path    string
	Pattern string
}

// MoveEntry describes a filesystem rename or move operation.
type MoveEntry struct {
	Src  string
	Dest string
}

// ContentReplaceEntry describes an in-place content replacement request.
type ContentReplaceEntry struct {
	Path       string
	OldContent string
	NewContent string
}

// SetPermissionEntry describes ownership and mode changes for a path.
type SetPermissionEntry struct {
	Path  string
	Mode  int
	Owner string
	Group string
}

// ReadFileOptions customizes file download behavior.
type ReadFileOptions struct {
	Encoding string
	Range    string
}
