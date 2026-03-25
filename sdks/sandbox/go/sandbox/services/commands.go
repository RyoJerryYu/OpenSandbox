package services

import (
	"context"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type ExecdCommands interface {
	Run(ctx context.Context, command string, opts *models.RunCommandOptions) (*models.CommandExecution, error)
	RunStream(ctx context.Context, command string, opts *models.RunCommandOptions) (*models.CommandStream, error)
	Interrupt(ctx context.Context, sessionID string) error
	GetCommandStatus(ctx context.Context, commandID string) (*models.CommandStatus, error)
	GetBackgroundCommandLogs(ctx context.Context, commandID string, cursor *int64) (*models.CommandLogs, error)
}
