package models

import "time"

// RunCommandOptions customizes command execution within a sandbox.
type RunCommandOptions struct {
	Background       bool
	WorkingDirectory string
	Timeout          *time.Duration
	UID              *int32
	GID              *int32
	Envs             map[string]string
}

// CommandStatus describes the state of a background command.
type CommandStatus struct {
	ID         string
	Command    string
	Running    bool
	ExitCode   *int32
	Error      string
	StartedAt  *time.Time
	FinishedAt *time.Time
}

// CommandLogs contains an incremental chunk of logs for a background command.
type CommandLogs struct {
	Content string
	Cursor  *int64
}

// OutputMessage represents one stdout or stderr line emitted during execution.
type OutputMessage struct {
	Text      string
	Timestamp int64
	IsError   bool
}

// ExecutionResult contains one logical result item returned by command/code execution.
type ExecutionResult struct {
	Text      string
	Timestamp int64
	Raw       map[string]any
}

// ExecutionError carries execution failure details.
type ExecutionError struct {
	Name      string
	Value     string
	Timestamp int64
	Traceback []string
}

// ExecutionLogs groups stdout and stderr messages collected during execution.
type ExecutionLogs struct {
	Stdout []OutputMessage
	Stderr []OutputMessage
}

// ExecutionComplete records completion timing for an execution.
type ExecutionComplete struct {
	Timestamp       int64
	ExecutionTimeMs int64
}

// ServerStreamEventError is the error payload embedded in streamed server events.
type ServerStreamEventError struct {
	Name      string   `json:"ename,omitempty"`
	Value     string   `json:"evalue,omitempty"`
	Traceback []string `json:"traceback,omitempty"`
}

// ServerStreamEvent is one event emitted by the streaming execution endpoint.
type ServerStreamEvent struct {
	Type           string                  `json:"type"`
	Timestamp      int64                   `json:"timestamp,omitempty"`
	Text           string                  `json:"text,omitempty"`
	Results        map[string]any          `json:"results,omitempty"`
	Error          *ServerStreamEventError `json:"error,omitempty"`
	ExecutionCount *int                    `json:"execution_count,omitempty"`
	ExecutionTime  *int64                  `json:"execution_time,omitempty"`
}

// CommandExecution is the aggregated result returned by non-streaming command execution.
type CommandExecution struct {
	ID             string
	ExecutionCount *int
	Result         []ExecutionResult
	Error          *ExecutionError
	Logs           ExecutionLogs
	Complete       *ExecutionComplete
}

// CommandStream exposes streaming command events and terminal completion state.
type CommandStream struct {
	Events <-chan ServerStreamEvent
	Done   <-chan error
}
