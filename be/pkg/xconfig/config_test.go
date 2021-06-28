package xconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadEnv(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		var output string
		setEnv(t, "TEST_CONFIG", "foo")
		MustLoadEnv("TEST_CONFIG", &output)
		require.Equal(t, "foo", output)
	})
	t.Run("Bool", func(t *testing.T) {
		var output bool
		setEnv(t, "TEST_CONFIG", "1")
		MustLoadEnv("TEST_CONFIG", &output)
		require.Equal(t, true, output)

		setEnv(t, "TEST_CONFIG", "0")
		MustLoadEnv("TEST_CONFIG", &output)
		require.Equal(t, false, output)
	})
	t.Run("Int", func(t *testing.T) {
		var output int
		setEnv(t, "TEST_CONFIG", "-100")
		MustLoadEnv("TEST_CONFIG", &output)
		require.Equal(t, -100, output)
	})
	t.Run("Uint", func(t *testing.T) {
		var output uint
		setEnv(t, "TEST_CONFIG", "100")
		MustLoadEnv("TEST_CONFIG", &output)
		require.Equal(t, uint(100), output)
	})
}

func setEnv(t *testing.T, env string, value string) {
	err := os.Setenv(env, value)
	require.NoError(t, err)
}
