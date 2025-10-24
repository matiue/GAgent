package main

import (
	"fmt"
	"net/http"
	"MonitorService/api"
)

func main() {
	// Register all monitoring endpoints
	http.HandleFunc("/metrics", api.MetricsHandler)   // raw network stats
	http.HandleFunc("/bandwidth", api.BandwidthHandler)
	http.HandleFunc("/cpu", api.CPUHandler)
	http.HandleFunc("/mem", api.MemHandler)
	http.HandleFunc("/disk", api.DiskHandler)
	http.HandleFunc("/system", api.SystemHandler)

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error:", err)
	}
}
