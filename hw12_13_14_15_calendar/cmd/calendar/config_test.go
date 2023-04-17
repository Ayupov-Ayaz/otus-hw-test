package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_WithDefaultConfigs(t *testing.T) {
	cfg, err := NewConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 8080, cfg.HTTP.Port)
	require.Equal(t, "debug", cfg.Logger.Level)
}

func TestNewConfig_WithCustomEnvironments(t *testing.T) {
	envs := map[string]string{
		"CALENDAR_HTTP_PORT":    "8081",
		"CALENDAR_LOGGER_LEVEL": "info",
	}

	for k, v := range envs {
		require.NoError(t, os.Setenv(k, v))
	}

	require.Equal(t, "8081", os.Getenv("CALENDAR_HTTP_PORT"))

	cfg, err := NewConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 8081, cfg.HTTP.Port)
	require.Equal(t, "info", cfg.Logger.Level)
}
