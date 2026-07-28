package config

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "0123456789abcdef0123456789abcdef" // 32 bytes

func TestLoad(t *testing.T) {
	env := func(overrides map[string]string) map[string]string {
		m := map[string]string{
			"PORT":              "8080",
			"DATABASE_URL":      "postgres://db",
			"BASE_URL":          "http://x",
			"LOG_LEVEL":         "info",
			"JWT_SECRET":        testJWTSecret,
			"ACCESS_TOKEN_TTL":  "5m",
			"REFRESH_TOKEN_TTL": "168h",
			"COOKIE_SECURE":     "false",
		}
		for k, v := range overrides {
			m[k] = v
		}
		return m
	}

	want := Config{
		Port:            "8080",
		DatabaseURL:     "postgres://db",
		BaseURL:         "http://x",
		LogLevel:        slog.LevelInfo,
		JWTSecret:       testJWTSecret,
		AccessTokenTTL:  5 * time.Minute,
		RefreshTokenTTL: 168 * time.Hour,
		CookieSecure:    false,
	}

	tcs := map[string]struct {
		env     map[string]string
		wantErr bool
		want    Config
	}{
		"valid": {
			env:  env(nil),
			want: want,
		},
		"log level debug": {
			env: env(map[string]string{"LOG_LEVEL": "debug"}),
			want: func() Config {
				c := want
				c.LogLevel = slog.LevelDebug
				return c
			}(),
		},
		"cookie secure true": {
			env: env(map[string]string{"COOKIE_SECURE": "true"}),
			want: func() Config {
				c := want
				c.CookieSecure = true
				return c
			}(),
		},
		"missing PORT":              {env: env(map[string]string{"PORT": ""}), wantErr: true},
		"missing DATABASE_URL":      {env: env(map[string]string{"DATABASE_URL": ""}), wantErr: true},
		"missing BASE_URL":          {env: env(map[string]string{"BASE_URL": ""}), wantErr: true},
		"missing LOG_LEVEL":         {env: env(map[string]string{"LOG_LEVEL": ""}), wantErr: true},
		"invalid LOG_LEVEL":         {env: env(map[string]string{"LOG_LEVEL": "verbose"}), wantErr: true},
		"PORT not a number":         {env: env(map[string]string{"PORT": "abc"}), wantErr: true},
		"PORT out of range":         {env: env(map[string]string{"PORT": "70000"}), wantErr: true},
		"missing JWT_SECRET":        {env: env(map[string]string{"JWT_SECRET": ""}), wantErr: true},
		"short JWT_SECRET":          {env: env(map[string]string{"JWT_SECRET": "tooshort"}), wantErr: true},
		"missing ACCESS_TOKEN_TTL":  {env: env(map[string]string{"ACCESS_TOKEN_TTL": ""}), wantErr: true},
		"invalid ACCESS_TOKEN_TTL":  {env: env(map[string]string{"ACCESS_TOKEN_TTL": "abc"}), wantErr: true},
		"missing REFRESH_TOKEN_TTL": {env: env(map[string]string{"REFRESH_TOKEN_TTL": ""}), wantErr: true},
		"missing COOKIE_SECURE":     {env: env(map[string]string{"COOKIE_SECURE": ""}), wantErr: true},
		"invalid COOKIE_SECURE":     {env: env(map[string]string{"COOKIE_SECURE": "maybe"}), wantErr: true},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			got, err := Load()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
