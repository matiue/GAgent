package config

import "time"

type Config struct {
	CollectInterval  time.Duration
	BatchSize        int
	QueueFile        string
	GRPCServer       string
	DiskPath         string
	NetworkInterface string
}

func LoadConfig() *Config {
	return &Config{
		CollectInterval: 5 * time.Second,
		BatchSize:       10,
		QueueFile:       "metrics.db",
		GRPCServer:      "localhost:50051",
		DiskPath:        "/",
		NetworkInterface: "eth0",
	}
}
