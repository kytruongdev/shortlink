package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Config holds infra/runtime settings sourced from environment variables.
type Config struct {
	Port        string
	DatabaseURL string
	BaseURL     string
	LogLevel    slog.Level
}

// MustLoad loads and validates config, panicking on error. Call only from main.
func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err.Error())
	}
	return cfg
}

// Load reads all config from env; every variable is required, no silent defaults.
func Load() (Config, error) {
	port, err := requireEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	databaseURL, err := requireEnv("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}

	baseURL, err := requireEnv("BASE_URL")
	if err != nil {
		return Config{}, err
	}

	logLevelStr, err := requireEnv("LOG_LEVEL")
	if err != nil {
		return Config{}, err
	}
	logLevel, err := parseLogLevel(logLevelStr)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:        port,
		DatabaseURL: databaseURL,
		BaseURL:     baseURL,
		LogLevel:    logLevel,
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks semantic constraints on the loaded configuration.
func (c Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("config: DATABASE_URL is required")
	}
	if p, err := strconv.Atoi(c.Port); err != nil || p < 1 || p > 65535 {
		return fmt.Errorf("config: PORT %q is not a valid port (1-65535)", c.Port)
	}
	return nil
}

// requireEnv returns the value of key, or an error when it is unset or empty.
func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("config: required env var %q is not set", key)
	}
	return v, nil
}

// parseLogLevel maps a level name to slog.Level, erroring on an unknown value.
func parseLogLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("config: invalid LOG_LEVEL %q (debug|info|warn|error)", s)
	}
}
