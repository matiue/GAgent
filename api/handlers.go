package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"Gaze/metrics"
)

// ---- Internal state to compute deltas ----
var (
	lock sync.Mutex

	prevNetStats   map[string]metrics.NetCounters
	prevNetSampleT time.Time

	prevCPUStats   metrics.CPUStats
	prevCPUSampleT time.Time
)

// ---------- NETWORK ----------
func MetricsHandler(w http.ResponseWriter, _ *http.Request) {
	stats, err := metrics.ParseProcNetDev()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func BandwidthHandler(w http.ResponseWriter, _ *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	curr, err := metrics.ParseProcNetDev()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now()
	if prevNetStats == nil {
		prevNetStats = curr
		prevNetSampleT = now
		json.NewEncoder(w).Encode(map[string]string{"message": "initial sample, retry later"})
		return
	}

	bw := metrics.MeasureBandwidth(prevNetStats, curr, now.Sub(prevNetSampleT).Seconds())
	prevNetStats = curr
	prevNetSampleT = now

	json.NewEncoder(w).Encode(bw)
}

// ---------- CPU ----------
func CPUHandler(w http.ResponseWriter, _ *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	curr, err := metrics.ParseProcStat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	now := time.Now()

	if prevCPUSampleT.IsZero() {
		prevCPUStats = curr
		prevCPUSampleT = now
		json.NewEncoder(w).Encode(map[string]string{"message": "initial CPU sample, retry later"})
		return
	}

	usage := metrics.CPUUsageFromDeltas(prevCPUStats, curr)
	prevCPUStats = curr
	prevCPUSampleT = now
	json.NewEncoder(w).Encode(usage)
}

// ---------- MEMORY ----------
func MemHandler(w http.ResponseWriter, _ *http.Request) {
	mem, err := metrics.ParseProcMeminfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(mem)
}

// ---------- DISK ----------
func DiskHandler(w http.ResponseWriter, _ *http.Request) {
	disks, err := metrics.GetAllDisks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(disks)
}

// ---------- SYSTEM SUMMARY ----------
func SystemHandler(w http.ResponseWriter, _ *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	summary := metrics.BuildSystemSummary(&prevCPUStats, &prevCPUSampleT, &prevNetStats, &prevNetSampleT)
	json.NewEncoder(w).Encode(summary)
}
