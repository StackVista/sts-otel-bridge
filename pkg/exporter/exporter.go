package exporter

import (
	"context"
	"os"
	"syscall"

	"github.com/stackvista/sts-otel-bridge/internal/config"
	"github.com/stackvista/sts-otel-bridge/internal/logging"
	"github.com/stackvista/sts-otel-bridge/pkg/batcher"
	"github.com/stackvista/sts-otel-bridge/pkg/collector"
	"github.com/stackvista/sts-otel-bridge/pkg/sender/stdout"
	"github.com/ztrue/shutdown"
	"google.golang.org/grpc"
)

type Exporter struct {
	cfg *config.Config
}

func NewExporter(cfg *config.Config) *Exporter {
	return &Exporter{
		cfg: cfg,
	}
}

func (e *Exporter) Run(ctx context.Context) error {
	server := grpc.NewServer()

	hooks := []func(os.Signal){}
	if e.cfg.Trace.Enabled {
		hook := startTrace(ctx, e.cfg.Trace, server)
		hooks = append(hooks, hook)
	}

	shutdown.AddWithParam(func(s os.Signal) {
		logger := logging.LoggerFor(ctx, "main-shutdown-hook")
		logger.Warn().Str("signal", s.String()).Msg("Received signal, shutting down GRPC server")
		server.Stop()

		logger.Warn().Msg("Shutting down all OTel subsystems")
		for _, hook := range hooks {
			hook(s)
		}
	})

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return nil
}

func startTrace(ctx context.Context, cfg config.Trace, server *grpc.Server) func(os.Signal) {
	tracer := collector.NewTraceCollector(cfg)
	traceBatcher := batcher.NewBatcher(cfg.Batcher, tracer.Out)
	sender := stdout.NewStdOutSender(traceBatcher.Out)
	tracer.Register(server)

	shutdownHook := func(s os.Signal) {
		logger := logging.LoggerFor(ctx, "tracer-shutdown-hook")
		logger.Warn().Str("signal", s.String()).Msg("Received signal, shutting down")
		tracer.Stop(ctx)

		traceBatcher.WaitGroup.Wait()
		sender.WaitGroup.Wait()
	}

	if err := sender.Start(ctx); err != nil {
		panic(err)
	}

	if err := traceBatcher.Start(ctx); err != nil {
		panic(err)
	}

	return shutdownHook
}
