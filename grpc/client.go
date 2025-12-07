package grpc

import (
	"context"
	"log"
	"time"

	pb "github.com/matiue/GAgent/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MetricsServiceClient
}

func NewClient(server string) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	client := pb.NewMetricsServiceClient(conn)
	return &Client{conn: conn, client: client}
}

func (c *Client) PushBatch(batch []*pb.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.PushMetrics(ctx, &pb.MetricsBatch{Metrics: batch})
	if err != nil {
		log.Printf("gRPC push error: %v", err)
		return
	}
	if resp.Success {
		log.Printf("Batch pushed successfully, %d metrics", len(batch))
	}
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
