package services

import (
	"context"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type ExecdCommands interface {
	Interrupt(ctx context.Context, sessionID string) error
	GetCommandStatus(ctx context.Context, commandID string) (*models.CommandStatus, error)
	GetBackgroundCommandLogs(ctx context.Context, commandID string, cursor *int64) (*models.CommandLogs, error)
}
