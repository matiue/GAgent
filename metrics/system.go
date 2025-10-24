package metrics

import (
	"time"
)

type SystemSummary struct {
	CPUUsage   *CPUUsage              `json:"cpu_usage,omitempty"`
	Memory     *MemInfo               `json:"memory,omitempty"`
	Bandwidth  map[string]*Bandwidth  `json:"bandwidth,omitempty"`
	Disks      []DiskUsage            `json:"disks,omitempty"`
	LastUpdate time.Time              `json:"last_update"`
}

// BuildSystemSummary returns a combined snapshot of system metrics.
func BuildSystemSummary(prevCPU *CPUStats, prevCPUT *time.Time, prevNet *map[string]NetCounters, prevNetT *time.Time) SystemSummary {
	summary := SystemSummary{LastUpdate: time.Now()}

	if prevCPU != nil && !prevCPUT.IsZero() {
		currCPU, _ := ParseProcStat()
		usage := CPUUsageFromDeltas(*prevCPU, currCPU)
		summary.CPUUsage = &usage
		*prevCPU = currCPU
		*prevCPUT = time.Now()
	}

	mem, _ := ParseProcMeminfo()
	summary.Memory = mem

	if prevNet != nil && !prevNetT.IsZero() {
		currNet, _ := ParseProcNetDev()
		bw := MeasureBandwidth(*prevNet, currNet, time.Since(*prevNetT).Seconds())
		summary.Bandwidth = bw
		*prevNet = currNet
		*prevNetT = time.Now()
	}

	disks, _ := GetAllDisks()
	summary.Disks = disks

	return summary
}
