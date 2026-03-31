package models

// SandboxMetrics contains a point-in-time resource usage snapshot reported by execd.
type SandboxMetrics struct {
	CPUCount          float32
	CPUUsedPercentage float32
	MemoryTotalMiB    float32
	MemoryUsedMiB     float32
	Timestamp         int64
}
