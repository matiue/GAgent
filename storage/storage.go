package storage

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	pb "github.com/matiue/GAgent/grpc/proto"
	grpcclient "github.com/matiue/GAgent/grpc"
)

type Queue struct {
	queue     []map[string]float64
	batchSize int
	mutex     sync.Mutex
	client    *grpcclient.Client
}

func NewQueue(queueFile string, batchSize int, client *grpcclient.Client) *Queue {
	return &Queue{
		queue:     make([]map[string]float64, 0),
		batchSize: batchSize,
		client:    client,
	}
}

func (q *Queue) Add(metric map[string]float64) {
	q.mutex.Lock()
	q.queue = append(q.queue, metric)
	if len(q.queue) >= q.batchSize {
		batch := q.queue
		q.queue = make([]map[string]float64, 0)
		q.mutex.Unlock()
		q.sendBatch(batch)
		return
	}
	q.mutex.Unlock()
}

func (q *Queue) Flush() {
	q.mutex.Lock()
	if len(q.queue) == 0 {
		q.mutex.Unlock()
		return
	}
	batch := q.queue
	q.queue = make([]map[string]float64, 0)
	q.mutex.Unlock()
	q.sendBatch(batch)
}

func (q *Queue) sendBatch(batch []map[string]float64) {
	if len(batch) == 0 {
		return
	}
	if q.client != nil {
		// Flatten maps to a slice of protobuf metrics
		var metrics []*pb.Metric
		for _, m := range batch {
			ts := time.Now().Unix()
			for name, value := range m {
				metrics = append(metrics, &pb.Metric{
					Name:      name,
					Value:     value,
					Timestamp: ts,
				})
			}
		}
		q.client.PushBatch(metrics)
		return
	}
	// Fallback: log JSON when no gRPC client is configured
	b, _ := json.Marshal(batch)
	log.Printf("Flushing metrics batch: %s\n", string(b))
}
