package config

import "time"

type RedisConfig struct {
	URL      string
	Password string
	DB       int
	Timeout  time.Duration
}
