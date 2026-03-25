package models

import "time"

type RunCommandOptions struct {
	Background       bool
	WorkingDirectory string
	Timeout          *time.Duration
	UID              *int32
	GID              *int32
	Envs             map[string]string
}

type CommandStatus struct {
	ID         string
	Command    string
	Running    bool
	ExitCode   *int32
	Error      string
	StartedAt  *time.Time
	FinishedAt *time.Time
}

type CommandLogs struct {
	Content string
	Cursor  *int64
}
