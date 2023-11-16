package config

import (
	"github.com/stackvista/sts-otel-bridge/pkg/batcher"
	"github.com/stackvista/sts-otel-bridge/pkg/grpc"
)

type Config struct {
	Trace Trace       `yaml:"trace" viper:"trace"`
	Grpc  grpc.Config `yaml:"grpc" viper:"grpc"`
}

type Trace struct {
	Enabled bool           `yaml:"enabled" default:"true"`
	Batcher batcher.Config `yaml:"batcher" viper:"batcher"`
}
