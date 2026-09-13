package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds infra/runtime settings sourced from environment variables.
type Config struct {
	Port            string
	DatabaseURL     string
	BaseURL         string
	LogLevel        slog.Level
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	CookieSecure    bool
	AllowedOrigins  []string
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

	jwtSecret, err := requireEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	accessTokenTTL, err := requireDurationEnv("ACCESS_TOKEN_TTL")
	if err != nil {
		return Config{}, err
	}

	refreshTokenTTL, err := requireDurationEnv("REFRESH_TOKEN_TTL")
	if err != nil {
		return Config{}, err
	}

	cookieSecure, err := requireBoolEnv("COOKIE_SECURE")
	if err != nil {
		return Config{}, err
	}

	allowedOrigins, err := requireCSVEnv("ALLOWED_ORIGINS")
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:            port,
		DatabaseURL:     databaseURL,
		BaseURL:         baseURL,
		LogLevel:        logLevel,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		CookieSecure:    cookieSecure,
		AllowedOrigins:  allowedOrigins,
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
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("config: JWT_SECRET must be at least 32 bytes")
	}
	if c.AccessTokenTTL <= 0 {
		return fmt.Errorf("config: ACCESS_TOKEN_TTL must be positive")
	}
	if c.RefreshTokenTTL <= 0 {
		return fmt.Errorf("config: REFRESH_TOKEN_TTL must be positive")
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

// requireDurationEnv reads key and parses it as a Go duration (e.g. "15m", "168h").
func requireDurationEnv(key string) (time.Duration, error) {
	v, err := requireEnv(key)
	if err != nil {
		return 0, err
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s %q is not a valid duration", key, v)
	}
	return d, nil
}

// requireCSVEnv reads key and splits it into a non-empty list of trimmed values.
func requireCSVEnv(key string) ([]string, error) {
	v, err := requireEnv(key)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("config: %s must list at least one value", key)
	}
	return out, nil
}

// requireBoolEnv reads key and parses it as a bool (true/false/1/0).
func requireBoolEnv(key string) (bool, error) {
	v, err := requireEnv(key)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("config: %s %q is not a valid bool", key, v)
	}
	return b, nil
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
