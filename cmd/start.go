package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/stackvista/sts-otel-bridge/internal/config"
	"github.com/stackvista/sts-otel-bridge/pkg/exporter"
)

func NewStartCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the sts-otel-bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context(), cfg)
		},
	}
}

func start(ctx context.Context, cfg *config.Config) error {
	exporter := exporter.NewExporter(cfg)

	return exporter.Run(ctx)
}
