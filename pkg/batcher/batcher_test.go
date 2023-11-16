package batcher

import (
	"context"
	"testing"
	"time"

	ststracepb "github.com/stackvista/sts-otel-bridge/proto/sts/trace"
	"github.com/stretchr/testify/assert"
)

func TestShouldBatch(t *testing.T) {
	ch := make(chan []*ststracepb.APITrace)
	batcher := NewBatcher(Config{2, time.Second, 1}, ch)
	assert.NotNil(t, batcher)
	err := batcher.Start(context.Background())
	assert.NoError(t, err)
	ch <- []*ststracepb.APITrace{
		{
			TraceID: 1,
		},
	}
	ch <- []*ststracepb.APITrace{
		{
			TraceID: 2,
		},
		{
			TraceID: 3,
		},
	}
	batch := <-batcher.Out
	assert.Equal(t, 2, len(batch))
	assert.Equal(t, uint64(1), batch[0].TraceID)
	assert.Equal(t, uint64(2), batch[1].TraceID)
	close(ch)
	batch = <-batcher.Out
	assert.Equal(t, 1, len(batch))
	assert.Equal(t, uint64(3), batch[0].TraceID)
}

func TestBatchTimeout(t *testing.T) {
	ch := make(chan []*ststracepb.APITrace)
	batcher := NewBatcher(Config{2, time.Second, 1}, ch)
	assert.NotNil(t, batcher)
	err := batcher.Start(context.Background())
	assert.NoError(t, err)
	ch <- []*ststracepb.APITrace{
		{
			TraceID: 1,
		},
	}
	time.Sleep(2 * time.Second) // Ensure 1 second timeout has passed
	batch := <-batcher.Out
	assert.Equal(t, 1, len(batch))
	assert.Equal(t, uint64(1), batch[0].TraceID)
}
