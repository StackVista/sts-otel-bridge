package stdout

import (
	"context"
	"sync"

	"github.com/stackvista/sts-otel-bridge/internal/logging"
	"github.com/stackvista/sts-otel-bridge/pkg/otel"
)

type StdOutSender[OTelData otel.OpenTelemetryData] struct {
	Name      string
	In        <-chan []OTelData
	WaitGroup *sync.WaitGroup
}

func NewStdOutSender[OTelData otel.OpenTelemetryData](in <-chan []OTelData, name string) *StdOutSender[OTelData] {
	return &StdOutSender[OTelData]{
		Name:      name,
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
	logger := logging.LoggerFor(ctx, s.Name)
	defer s.WaitGroup.Done()

	for {
		select {
		case d, ok := <-s.In:
			if !ok {
				logger.Warn().Msg("Channel closed, exiting")
				return
			}

			logger.Debug().Int("data-length", len(d)).Msg("Received data")

			for _, item := range d {
				otel.PrintOTelData(item)
			}
		}
	}
}
