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

type OutputMessage struct {
	Text      string
	Timestamp int64
	IsError   bool
}

type ExecutionResult struct {
	Text      string
	Timestamp int64
	Raw       map[string]any
}

type ExecutionError struct {
	Name      string
	Value     string
	Timestamp int64
	Traceback []string
}

type ExecutionLogs struct {
	Stdout []OutputMessage
	Stderr []OutputMessage
}

type ExecutionComplete struct {
	Timestamp       int64
	ExecutionTimeMs int64
}

type ServerStreamEventError struct {
	Name      string   `json:"ename,omitempty"`
	Value     string   `json:"evalue,omitempty"`
	Traceback []string `json:"traceback,omitempty"`
}

type ServerStreamEvent struct {
	Type           string                  `json:"type"`
	Timestamp      int64                   `json:"timestamp,omitempty"`
	Text           string                  `json:"text,omitempty"`
	Results        map[string]any          `json:"results,omitempty"`
	Error          *ServerStreamEventError `json:"error,omitempty"`
	ExecutionCount *int                    `json:"execution_count,omitempty"`
	ExecutionTime  *int64                  `json:"execution_time,omitempty"`
}

type CommandExecution struct {
	ID             string
	ExecutionCount *int
	Result         []ExecutionResult
	Error          *ExecutionError
	Logs           ExecutionLogs
	Complete       *ExecutionComplete
}

type CommandStream struct {
	Events <-chan ServerStreamEvent
	Done   <-chan error
}
