package utils

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// CPU usage as % (simplified)
func ReadCPUUsage() float64 {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				return 0
			}
			var total float64
			var idle float64
			for i := 1; i < len(fields); i++ {
				v, err := strconv.ParseFloat(fields[i], 64)
				if err != nil {
					v = 0
				}
				total += v
				if i == 4 { // idle is the 4th field (index 4)
					idle = v
				}
			}
			if total == 0 {
				return 0
			}
			return (total - idle) / total * 100
		}
	}
	return 0
}

// Memory usage as % (simplified)
func ReadMemoryUsage() float64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	lines := strings.Split(string(data), "\n")
	var memTotal, memAvailable float64
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "MemTotal:" {
			memTotal, _ = strconv.ParseFloat(fields[1], 64)
		} else if fields[0] == "MemAvailable:" {
			memAvailable, _ = strconv.ParseFloat(fields[1], 64)
		}
	}
	if memTotal == 0 {
		return 0
	}
	return (memTotal - memAvailable) / memTotal * 100
}

// Disk usage % (simplified)
func ReadDiskUsage(path string) float64 {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return 0
	}
	total := float64(stat.Blocks) * float64(stat.Bsize)
	free := float64(stat.Bfree) * float64(stat.Bsize)
	if total == 0 {
		return 0
	}
	return (total - free) / total * 100
}

// Network usage as bytes transmitted (simplified)
func ReadNetworkUsage(interfaceName string) float64 {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, interfaceName+":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			return 0
		}
		fields := strings.Fields(parts[1])
		if len(fields) < 9 {
			return 0
		}
		rx, _ := strconv.ParseFloat(fields[0], 64) // rx bytes
		tx, _ := strconv.ParseFloat(fields[8], 64) // tx bytes
		return rx + tx
	}
	return 0
}
