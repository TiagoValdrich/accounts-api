package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	EnvVars EnvironmentVariables
}

type EnvironmentVariables struct {
	Database DatabaseConfig
}

func load() (*AppConfig, error) {
	cfg := &AppConfig{}

	setupLogger()

	if err := env.Parse(&cfg.EnvVars); err != nil {
		log.Err(err).Msg("failed to load environment variables")

		return nil, err
	}
	return cfg, nil
}

func setupLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func MustLoad() *AppConfig {
	cfg, err := load()
	if err != nil {
		panic(err)
	}
	return cfg
}
