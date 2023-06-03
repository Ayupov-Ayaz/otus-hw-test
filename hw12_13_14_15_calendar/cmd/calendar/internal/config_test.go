package internal

import (
	"os"
	"testing"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/parser"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalEnv_WithDefaultConfigs(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		envs    map[string]string
		checker func(t *testing.T, cfg *Config)
	}{
		{
			name: "with default configs",
			err:  nil,
			checker: func(t *testing.T, cfg *Config) {
				require.Equal(t, 8080, cfg.HTTP.Port)
				require.Equal(t, "debug", cfg.Logger.Level)
			},
		},
		{
			name: "unmarshal with custom envs",
			err:  nil,
			envs: map[string]string{
				"CALENDAR_HTTP_PORT":    "8081",
				"CALENDAR_LOGGER_LEVEL": "info",
			},
			checker: func(t *testing.T, cfg *Config) {
				require.Equal(t, 8081, cfg.HTTP.Port)
				require.Equal(t, "info", cfg.Logger.Level)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				require.NoError(t, os.Setenv(k, v))
			}

			cfg := &Config{}
			err := parser.UnmarshalEnv(envPrefix, cfg)
			require.ErrorIs(t, err, tt.err)
			require.NotNil(t, cfg)
			tt.checker(t, cfg)
		})
	}
}

func createYamlFile(t *testing.T, data string) *os.File {
	t.Helper()
	f, err := os.CreateTemp("./", "config.*.yaml")
	require.NoError(t, err)

	_, err = f.WriteString(data)
	require.NoError(t, err)

	return f
}

func TestUnmarshalYaml_WithCustomConfigs(t *testing.T) {
	const data = `http:
  port: 8081
logger:
  level: info`

	f := createYamlFile(t, data)
	err := os.Remove(f.Name())
	require.NoError(t, err)

	cfg := &Config{}
	err = parser.UnmarshalYaml([]byte(data), cfg)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 8081, cfg.HTTP.Port)
	require.Equal(t, "info", cfg.Logger.Level)
}
