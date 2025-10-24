package metrics

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type CPUStats struct {
	User, Nice, System, Idle, IOWait, IRQ, SoftIRQ, Steal uint64
}

type CPUUsage struct {
	UserPercent   float64 `json:"user_percent"`
	SystemPercent float64 `json:"system_percent"`
	IdlePercent   float64 `json:"idle_percent"`
	TotalPercent  float64 `json:"total_percent"`
}

// ParseProcStat reads /proc/stat for CPU time counters.
func ParseProcStat() (CPUStats, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return CPUStats{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 7 && fields[0] == "cpu" {
			toUint := func(s string) uint64 {
				v, _ := strconv.ParseUint(s, 10, 64)
				return v
			}
			return CPUStats{
				User:    toUint(fields[1]),
				Nice:    toUint(fields[2]),
				System:  toUint(fields[3]),
				Idle:    toUint(fields[4]),
				IOWait:  toUint(fields[5]),
				IRQ:     toUint(fields[6]),
				SoftIRQ: toUint(fields[7]),
				Steal:   func() uint64 { if len(fields) > 8 { return toUint(fields[8]) }; return 0 }(),
			}, nil
		}
	}
	return CPUStats{}, scanner.Err()
}

// CPUUsageFromDeltas calculates usage percentages.
func CPUUsageFromDeltas(prev, curr CPUStats) CPUUsage {
	prevTotal := prev.User + prev.Nice + prev.System + prev.Idle + prev.IOWait + prev.IRQ + prev.SoftIRQ + prev.Steal
	currTotal := curr.User + curr.Nice + curr.System + curr.Idle + curr.IOWait + curr.IRQ + curr.SoftIRQ + curr.Steal
	totalDiff := float64(currTotal - prevTotal)
	idleDiff := float64(curr.Idle - prev.Idle)
	systemDiff := float64(curr.System - prev.System)
	userDiff := float64(curr.User - prev.User)

	var usage CPUUsage
	if totalDiff > 0 {
		usage.IdlePercent = (idleDiff / totalDiff) * 100
		usage.SystemPercent = (systemDiff / totalDiff) * 100
		usage.UserPercent = (userDiff / totalDiff) * 100
		usage.TotalPercent = 100 - usage.IdlePercent
	}
	return usage
}
