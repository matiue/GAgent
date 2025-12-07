package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matiue/GAgent/collector"
	"github.com/matiue/GAgent/config"
	grpcclient "github.com/matiue/GAgent/grpc"
	"github.com/matiue/GAgent/storage"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize gRPC client
	client := grpcclient.NewClient(cfg.GRPCServer)
	defer client.Close()

	// Initialize storage queue
	storageQueue := storage.NewQueue(cfg.QueueFile, cfg.BatchSize, client)

	// Start collectors
	ticker := time.NewTicker(cfg.CollectInterval)
	defer ticker.Stop()

	// Handle termination signals to flush remaining metrics
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			metrics := make(map[string]float64)
			metrics["cpu"] = collector.GetCPUUsage()
			metrics["memory"] = collector.GetMemoryUsage()
			metrics["disk"] = collector.GetDiskUsage(cfg.DiskPath)
			metrics["network"] = collector.GetNetworkUsage(cfg.NetworkInterface)

			// Push to local queue
			storageQueue.Add(metrics)
		case sig := <-sigCh:
			log.Printf("Received signal %v, flushing queue and exiting", sig)
			storageQueue.Flush()
			return
		}
	}
}
