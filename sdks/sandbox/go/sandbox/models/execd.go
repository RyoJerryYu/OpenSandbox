package models

type SandboxMetrics struct {
	CPUCount          float32
	CPUUsedPercentage float32
	MemoryTotalMiB    float32
	MemoryUsedMiB     float32
	Timestamp         int64
}
