package models

import "time"

type EntryInfo struct {
	Path       string
	Size       int64
	ModifiedAt time.Time
	CreatedAt  time.Time
	Mode       int
	Owner      string
	Group      string
}

type WriteEntry struct {
	Path     string
	Data     any
	Mode     int
	Owner    string
	Group    string
	Encoding string
}

type SearchEntry struct {
	Path    string
	Pattern string
}

type MoveEntry struct {
	Src  string
	Dest string
}

type ContentReplaceEntry struct {
	Path       string
	OldContent string
	NewContent string
}

type SetPermissionEntry struct {
	Path  string
	Mode  int
	Owner string
	Group string
}

type ReadFileOptions struct {
	Encoding string
	Range    string
}
