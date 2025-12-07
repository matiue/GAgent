package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/matiue/GAgent/grpc/proto"
	"google.golang.org/grpc"
)

// Server implements the MetricsService gRPC interface
type server struct {
	pb.UnimplementedMetricsServiceServer
}

// PushMetrics handles incoming metric batches
func (s *server) PushMetrics(ctx context.Context, batch *pb.MetricsBatch) (*pb.PushResponse, error) {
	fmt.Printf("\nReceived batch at %s:\n", time.Now().Format(time.RFC3339))
	for _, m := range batch.Metrics {
		fmt.Printf(" - Metric: %s, Value: %.2f, Timestamp: %d\n", m.Name, m.Value, m.Timestamp)
	}
	return &pb.PushResponse{Success: true}, nil
}

func main() {
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterMetricsServiceServer(s, &server{})

	log.Println("gRPC server running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
