package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hierynomus/autobind"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackvista/sts-otel-bridge/internal/config"
)

const (
	VerboseFlag      = "verbose"
	VerboseFlagShort = "v"
)

func NewRootCmd(cfg *config.Config) *cobra.Command {
	var verbosity int

	cmd := &cobra.Command{
		Use:   "sts-otel-bridge",
		Short: "sts-otel-bridge is a collector for OpenTelemetry traces and metrics",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch verbosity {
			case 0:
				// Nothing to do
			case 1:
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			default:
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			}

			vp := viper.New()
			vp.SetConfigName("config")
			vp.AddConfigPath(".")
			vp.SetConfigType("yaml")

			logger := log.Ctx(cmd.Context())

			if err := vp.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					logger.Warn().Msg("No config file found... Continuing with defaults")
					// Config file not found; ignore error if desired
				} else {
					fmt.Printf("%s", err)
					os.Exit(1)
				}
			}

			binder := &autobind.Autobinder{
				UseNesting:   true,
				EnvPrefix:    "K2P",
				ConfigObject: cfg,
				Viper:        vp,
				SetDefaults:  true,
			}

			binder.Bind(cmd.Context(), cmd, []string{})

			return nil
		},
	}
	cmd.PersistentFlags().CountVarP(&verbosity, VerboseFlag, VerboseFlagShort, "Print verbose logging to the terminal (use multiple times to increase verbosity)")

	return cmd
}

func Execute(ctx context.Context) {
	cfg := &config.Config{}

	rootCmd := NewRootCmd(cfg)
	rootCmd.AddCommand(NewStartCmd(cfg))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
