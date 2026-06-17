package config

import (
	"flag"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Http        HttpConfig
	Db          DbConfig
	Argon2id    Argon2idConfig
	Jwt         JwtConfig
	GoogleOAuth GoogleOAuthConfig
	Redis       RedisConfig
}

type HttpConfig struct {
	Port string `env:"HTTP_PORT" env-required:"true"`
}

type DbConfig struct {
	Username string `env:"DB_USERNAME" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
	SSLMode  string `env:"DB_SSL_MODE" env-required:"true"`
}

type Argon2idConfig struct {
	Iteration uint32 `env:"ARGON2ID_ITERATION" env-required:"true"`
	MemoryMB  uint32 `env:"ARGON2ID_MEMORY_MB" env-required:"true"`
	Threads   uint8  `env:"ARGON2ID_THREADS" env-required:"true"`
	Key       uint32 `env:"ARGON2ID_KEY" env-required:"true"`
}

type JwtConfig struct {
	AccessSecret         string `env:"JWT_SECRET_ACCESS" env-required:"true"`
	RefreshSecret        string `env:"JWT_SECRET_REFRESH" env-required:"true"`
	AccessDurationInMin  int    `env:"JWT_ACCESS_DURATION_IN_MIN" env-required:"true"`
	RefreshDurationInDay int    `env:"JWT_REFRESH_DURATION_DAY" env-required:"true"`
}

type GoogleOAuthConfig struct {
	WebClientID string `env:"WEB_CLIENT_ID" env-required:"true"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     string `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `env:"REDIS_DB" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string
	var cfg *Config

	cfg = &Config{}

	err := cleanenv.ReadEnv(cfg)

	if err == nil {
		return cfg
	}

	log.Info().Msg("failed to read env, falling back to CONFIG_PATH")

	configPath = strings.TrimSpace(os.Getenv("CONFIG_PATH"))

	if configPath == "" {
		flags := flag.String("config", ".env", "config file path")
		flag.Parse()
		configPath = *flags
	}

	if configPath == "" {
		log.Info().Msg("config file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Msg("config file does not exist")
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatal().Err(err).Msg("error reading config")
	}

	return cfg
}
