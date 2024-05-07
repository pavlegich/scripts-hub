package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	want := &Config{}
	t.Run("success", func(t *testing.T) {
		got := NewConfig(context.Background())
		require.Equal(t, want, got)
	})
}

func TestConfig_ParseFlags(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.ParseFlags(context.Background())
		require.NoError(t, err)
	})
}
