package metrics

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type MemInfo struct {
	MemTotalKB     uint64  `json:"mem_total_kb"`
	MemAvailableKB uint64  `json:"mem_available_kb"`
	MemUsedKB      uint64  `json:"mem_used_kb"`
	UsedPercent    float64 `json:"used_percent"`

	SwapTotalKB     uint64  `json:"swap_total_kb"`
	SwapFreeKB      uint64  `json:"swap_free_kb"`
	SwapUsedKB      uint64  `json:"swap_used_kb"`
	SwapUsedPercent float64 `json:"swap_used_percent"`

	DirtyKB        uint64 `json:"dirty_kb"`
	WritebackKB    uint64 `json:"writeback_kb"`
	CachedKB       uint64 `json:"cached_kb"`
	BuffersKB      uint64 `json:"buffers_kb"`
	SReclaimableKB uint64 `json:"sreclaimable_kb"`

	ZswapKB    uint64 `json:"zswap_kb"`    // compressed pages in RAM
	ZswappedKB uint64 `json:"zswapped_kb"` // total data written to zswap
}

// ParseProcMeminfo reads /proc/meminfo and fills MemInfo
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
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			// skip invalid entries
			continue
		}

		switch key {
		case "MemTotal":
			mem.MemTotalKB = val
		case "MemAvailable":
			mem.MemAvailableKB = val
		case "SwapTotal":
			mem.SwapTotalKB = val
		case "SwapFree":
			mem.SwapFreeKB = val
		case "Dirty":
			mem.DirtyKB = val
		case "Writeback":
			mem.WritebackKB = val
		case "Cached":
			mem.CachedKB = val
		case "Buffers":
			mem.BuffersKB = val
		case "SReclaimable":
			mem.SReclaimableKB = val
		case "Zswap":
			mem.ZswapKB = val
		case "Zswapped":
			mem.ZswappedKB = val
		}
	}

	// Compute derived metrics safely
	if mem.MemTotalKB > 0 && mem.MemAvailableKB > 0 {
		mem.MemUsedKB = mem.MemTotalKB - mem.MemAvailableKB
		mem.UsedPercent = float64(mem.MemUsedKB) / float64(mem.MemTotalKB) * 100
	}

	if mem.SwapTotalKB > 0 && mem.SwapFreeKB > 0 {
		mem.SwapUsedKB = mem.SwapTotalKB - mem.SwapFreeKB
		mem.SwapUsedPercent = float64(mem.SwapUsedKB) / float64(mem.SwapTotalKB) * 100
	}

	return mem, scanner.Err()
}
