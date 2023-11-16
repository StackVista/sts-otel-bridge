package stdout

import (
	"context"
	"sync"

	"github.com/stackvista/sts-otel-bridge/internal/logging"
	"github.com/stackvista/sts-otel-bridge/pkg/otel"
)

type StdOutSender[OTelData otel.OpenTelemetryData] struct {
	In        <-chan []OTelData
	WaitGroup *sync.WaitGroup
}

func NewStdOutSender[OTelData otel.OpenTelemetryData](in <-chan []OTelData) *StdOutSender[OTelData] {
	return &StdOutSender[OTelData]{
		In:        in,
		WaitGroup: &sync.WaitGroup{},
	}
}

func (s *StdOutSender[_]) Start(ctx context.Context) error {
	s.WaitGroup.Add(1)
	go s.run(ctx)

	return nil
}

func (s *StdOutSender[DataList]) run(ctx context.Context) {
	logger := logging.LoggerFor(ctx, "stdout-sender")
	defer s.WaitGroup.Done()

	for {
		select {
		case ts, ok := <-s.In:
			if !ok {
				logger.Warn().Msg("Channel closed, exiting")
				return
			}

			for _, t := range ts {
				otel.PrintOTelData(t)
			}
		}
	}
}
