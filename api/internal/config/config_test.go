package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	env := func(overrides map[string]string) map[string]string {
		m := map[string]string{
			"PORT":         "8080",
			"DATABASE_URL": "postgres://db",
			"BASE_URL":     "http://x",
			"LOG_LEVEL":    "info",
		}
		for k, v := range overrides {
			m[k] = v
		}
		return m
	}

	tcs := map[string]struct {
		env     map[string]string
		wantErr bool
		want    Config
	}{
		"valid": {
			env:  env(nil),
			want: Config{Port: "8080", DatabaseURL: "postgres://db", BaseURL: "http://x", LogLevel: slog.LevelInfo},
		},
		"log level debug": {
			env:  env(map[string]string{"LOG_LEVEL": "debug"}),
			want: Config{Port: "8080", DatabaseURL: "postgres://db", BaseURL: "http://x", LogLevel: slog.LevelDebug},
		},
		"missing PORT":         {env: env(map[string]string{"PORT": ""}), wantErr: true},
		"missing DATABASE_URL": {env: env(map[string]string{"DATABASE_URL": ""}), wantErr: true},
		"missing BASE_URL":     {env: env(map[string]string{"BASE_URL": ""}), wantErr: true},
		"missing LOG_LEVEL":    {env: env(map[string]string{"LOG_LEVEL": ""}), wantErr: true},
		"invalid LOG_LEVEL":    {env: env(map[string]string{"LOG_LEVEL": "verbose"}), wantErr: true},
		"PORT not a number":    {env: env(map[string]string{"PORT": "abc"}), wantErr: true},
		"PORT out of range":    {env: env(map[string]string{"PORT": "70000"}), wantErr: true},
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
