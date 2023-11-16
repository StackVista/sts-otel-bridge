package logging

import (
	"context"

	"github.com/rs/zerolog"
)

func LoggerFor(ctx context.Context, name string) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("component", name).Logger()
}
