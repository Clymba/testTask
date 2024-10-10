package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Name string `env:"APP_NAME"     envDefault:"AttemptOfCleanProject"`
	Port string `env:"ADDRESS"      envDefault:":8080"`
	DSN  string `env:"POSTGRES_DSN" envDefault:"postgres://postgres:postgres@localhost:5433/postgresDB"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
