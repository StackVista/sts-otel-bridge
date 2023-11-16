package exporter

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/stackvista/sts-otel-bridge/internal/config"
	"github.com/stackvista/sts-otel-bridge/internal/logging"
	"github.com/stackvista/sts-otel-bridge/pkg/batcher"
	"github.com/stackvista/sts-otel-bridge/pkg/collector"
	"github.com/stackvista/sts-otel-bridge/pkg/sender/stdout"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
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
	logger := logging.LoggerFor(ctx, "main")
	server := grpc.NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", e.cfg.Grpc.Port))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to listen")
	}

	hooks := []func(os.Signal){}
	if e.cfg.Trace.Enabled {
		hook := startTrace(ctx, e.cfg.Trace, server)
		hooks = append(hooks, hook)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := <-sigCh
		logger := logging.LoggerFor(ctx, "main-hook")
		logger.Warn().Str("signal", s.String()).Msg("Received signal, shutting down GRPC server")
		server.Stop()

		logger.Warn().Msg("Shutting down all OTel subsystems")
		for _, hook := range hooks {
			hook(s)
		}
		logger.Info().Msg("Shutdown complete")
	}()

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("could not serve: %v", err)
	}
	wg.Wait()
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
