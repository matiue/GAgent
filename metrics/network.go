package metrics

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"time"
)

type NetCounters struct {
	Name        string    `json:"name"`
	BytesRecv   uint64    `json:"bytes_recv"`
	BytesSent   uint64    `json:"bytes_sent"`
	PacketsRecv uint64    `json:"packets_recv"`
	PacketsSent uint64    `json:"packets_sent"`
	RecvErrs    uint64    `json:"recv_errs"`
	SendErrs    uint64    `json:"send_errs"`
	RecvDrops   uint64    `json:"recv_drops"`
	SendDrops   uint64    `json:"send_drops"`
	LastUpdate  time.Time `json:"LastUpdate"`
}

type Bandwidth struct {
	Name            string  `json:"name"`
	BytesRecvPerSec float64 `json:"bytes_recv_per_sec"`
	BytesSentPerSec float64 `json:"bytes_sent_per_sec"`
}

// ParseProcNetDev parses /proc/net/dev to get raw interface counters.
func ParseProcNetDev() (map[string]NetCounters, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := make(map[string]NetCounters)
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if i < 2 || line == "" {
			continue
		}
		fields := strings.Fields(strings.ReplaceAll(line, ":", " "))
		if len(fields) < 17 {
			continue
		}

		name := fields[0]
		var nc NetCounters
		nc.Name = name
		fmtSscanf(fields[1:], &nc.BytesRecv, &nc.PacketsRecv, &nc.RecvErrs, &nc.RecvDrops,
			new(uint64), new(uint64), new(uint64), new(uint64), // placeholders
			&nc.BytesSent, &nc.PacketsSent, &nc.SendErrs, &nc.SendDrops)
		nc.LastUpdate = time.Now()
		stats[name] = nc
	}
	return stats, scanner.Err()
}

// MeasureBandwidth computes bytes/sec between two samples.
func MeasureBandwidth(prev, curr map[string]NetCounters, dt float64) map[string]*Bandwidth {
	result := make(map[string]*Bandwidth)
	for name, c := range curr {
		if p, ok := prev[name]; ok && dt > 0 {
			result[name] = &Bandwidth{
				Name:            name,
				BytesRecvPerSec: float64(c.BytesRecv-p.BytesRecv) / dt,
				BytesSentPerSec: float64(c.BytesSent-p.BytesSent) / dt,
			}
		}
	}
	return result
}

// fmtSscanf parses multiple uint64s from string fields.
func fmtSscanf(fields []string, vals ...*uint64) {
	for i := 0; i < len(fields) && i < len(vals); i++ {
		var v uint64
		_, _ = fmt.Sscan(fields[i], &v)
		*vals[i] = v
	}
}
