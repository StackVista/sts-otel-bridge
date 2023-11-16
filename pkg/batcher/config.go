package batcher

import "time"

type Config struct {
	BatchSize    int           `yaml:"batch-size" viper:"batch-size" env:"BATCH_SIZE" default:"100"`
	BatchTimeout time.Duration `yaml:"timeout" viper:"timeout" env:"TIMEOUT" default:"30s"`
	BufferSize   int           `yaml:"buffer-size" default:"0" viper:"buffer-size" env:"BUFFER_SIZE"`
}
