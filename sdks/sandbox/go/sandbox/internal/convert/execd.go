package convert

import (
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func FromExecdMetrics(resp *execdapi.Metrics) *models.SandboxMetrics {
	if resp == nil {
		return nil
	}
	return &models.SandboxMetrics{
		CPUCount:          resp.CpuCount,
		CPUUsedPercentage: resp.CpuUsedPct,
		MemoryTotalMiB:    resp.MemTotalMib,
		MemoryUsedMiB:     resp.MemUsedMib,
		Timestamp:         resp.Timestamp,
	}
}

func FromExecdCommandStatus(resp *execdapi.CommandStatusResponse) *models.CommandStatus {
	if resp == nil {
		return nil
	}

	status := &models.CommandStatus{}
	if resp.Id != nil {
		status.ID = *resp.Id
	}
	if resp.Content != nil {
		status.Command = *resp.Content
	}
	if resp.Running != nil {
		status.Running = *resp.Running
	}
	if resp.ExitCode != nil {
		exitCode := *resp.ExitCode
		status.ExitCode = &exitCode
	}
	if resp.Error != nil {
		status.Error = *resp.Error
	}
	if resp.StartedAt != nil {
		startedAt := *resp.StartedAt
		status.StartedAt = &startedAt
	}
	if resp.FinishedAt != nil {
		finishedAt := *resp.FinishedAt
		status.FinishedAt = &finishedAt
	}
	return status
}
