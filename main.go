//go:generate protoc --proto_path=proto --go_out=. --go_opt=module=github.com/stackvista/sts-otel-bridge proto/trace_payload.proto proto/trace.proto proto/span.proto
package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stackvista/sts-otel-bridge/cmd"
)

func main() {
	ctx := log.Logger.WithContext(context.Background())
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cmd.Execute(ctx)
}
