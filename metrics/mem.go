package metrics

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type MemInfo struct {
	MemTotalKB     uint64  `json:"mem_total_kb"`
	MemFreeKB      uint64  `json:"mem_free_kb"`
	MemAvailableKB uint64  `json:"mem_available_kb"`
	UsedPercent    float64 `json:"used_percent"`
}

// ParseProcMeminfo parses /proc/meminfo for RAM stats.
func ParseProcMeminfo() (*MemInfo, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mem := &MemInfo{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		switch key {
		case "MemTotal":
			mem.MemTotalKB = val
		case "MemFree":
			mem.MemFreeKB = val
		case "MemAvailable":
			mem.MemAvailableKB = val
		}
	}
	if mem.MemTotalKB > 0 {
		used := mem.MemTotalKB - mem.MemAvailableKB
		mem.UsedPercent = float64(used) / float64(mem.MemTotalKB) * 100
	}
	return mem, scanner.Err()
}
