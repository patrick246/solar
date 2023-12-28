package config

import (
	"log/slog"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel    slog.Level `envconfig:"LOG_LEVEL" default:"INFO"`
	MetricsAddr string     `envconfig:"METRICS_ADDRESS" default:":8080"`

	Database Database
	Broker   Broker
}

type Database struct {
	URL          string `envconfig:"DB_URL" default:"postgresql://postgres:postgres@localhost:5432/postgres"`
	MaxOpenConns int    `envconfig:"DB_MAX_OPEN_CONNS" default:"16"`
}

type Broker struct {
	URL   string `envconfig:"MQTT_URL" default:"mqtt://10.76.0.101:1883"`
	Topic string `envconfig:"MQTT_TOPIC" default:"shellypro3em-3ce90e7025c8/events/rpc"`
}

func Get() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
