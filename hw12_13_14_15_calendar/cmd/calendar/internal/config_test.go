package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalEnv_WithDefaultConfigs(t *testing.T) {
	cfg := &Config{}
	err := unmarshalEnv(cfg)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 8080, cfg.HTTP.Port)
	require.Equal(t, "debug", cfg.Logger.Level)
}

func TestUnmarshalEnv_WithCustomEnvironments(t *testing.T) {
	envs := map[string]string{
		"CALENDAR_HTTP_PORT":    "8081",
		"CALENDAR_LOGGER_LEVEL": "info",
	}

	for k, v := range envs {
		require.NoError(t, os.Setenv(k, v))
	}

	require.Equal(t, "8081", os.Getenv("CALENDAR_HTTP_PORT"))

	cfg := &Config{}
	err := unmarshalEnv(cfg)
	require.NoError(t, err)
	require.Equal(t, 8081, cfg.HTTP.Port)
	require.Equal(t, "info", cfg.Logger.Level)
}

func createYamlFile(t *testing.T, data string) *os.File {
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
	os.Remove(f.Name())

	cfg := &Config{}
	err := unmarshalYaml([]byte(data))(cfg)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 8081, cfg.HTTP.Port)
	require.Equal(t, "info", cfg.Logger.Level)
}
