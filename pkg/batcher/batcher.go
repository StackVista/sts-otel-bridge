package batcher

import (
	"context"
	"sync"
	"time"

	"github.com/stackvista/sts-otel-bridge/internal/logging"
	"github.com/stackvista/sts-otel-bridge/pkg/otel"
)

type Batcher[OTelData otel.OpenTelemetryData] struct {
	BatchSize    int
	BatchTimeout time.Duration
	In           <-chan []OTelData
	Out          chan []OTelData
	WaitGroup    *sync.WaitGroup
}

func NewBatcher[D otel.OpenTelemetryData](config Config, in <-chan []D) *Batcher[D] {
	return &Batcher[D]{
		BatchSize:    config.BatchSize,
		BatchTimeout: config.BatchTimeout,
		In:           in,
		Out:          make(chan []D, config.BufferSize),
		WaitGroup:    &sync.WaitGroup{},
	}
}

func (b *Batcher[_]) Start(ctx context.Context) error {
	b.WaitGroup.Add(1)

	go b.Run(ctx)

	return nil
}

func (b *Batcher[DataList]) Run(ctx context.Context) {
	logger := logging.LoggerFor(ctx, "batcher")
	batch := []DataList{}
	ticker := time.NewTicker(b.BatchTimeout)

	defer b.WaitGroup.Done()
	defer close(b.Out)

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Batcher stopping...")
			b.Out <- batch
			return
		case ts, ok := <-b.In:
			if !ok {
				logger.Info().Msg("Incoming channel closed, flush remaining batch")
				b.Out <- batch
				return
			}

			logger.Trace().Int("batch", len(batch)).Int("incoming", len(ts)).Msg("Received data...")
			batch = append(batch, ts...)
			if len(batch) >= b.BatchSize {
				logger.Debug().Int("batch", len(batch)).Msg("Batch full, sending...")
				b.Out <- batch[0:b.BatchSize]
				batch = batch[b.BatchSize:]
				ticker.Reset(b.BatchTimeout)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				logger.Debug().Int("batch", len(batch)).Msg("Batch timeout, sending...")
				b.Out <- batch
				batch = []DataList{}
			}
		}
	}
}
