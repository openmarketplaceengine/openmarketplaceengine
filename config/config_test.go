package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	override := "123456"

	err := os.Setenv(fmt.Sprintf("%s_%s", EnvPrefix, "SERVICE_PORT"), override)
	require.NoError(t, err)

	value := os.Getenv(fmt.Sprintf("%s_%s", EnvPrefix, "SERVICE_PORT"))
	require.Equal(t, override, value)

	err = Read()
	require.NoError(t, err)

	t.Run("testReadYml", func(t *testing.T) {
		testReadYml(t)
	})
	t.Run("testEnvVarOverride", func(t *testing.T) {
		testEnvVarOverride(t, override)
	})
	t.Run("testMissingPropertyError", func(t *testing.T) {
		testMissingPropertyError(t)
	})

	t.Run("testAllRequired", func(t *testing.T) {
		testAllRequired(t)
	})
}

func testReadYml(t *testing.T) {
	addr := GetString(RedisStoreAddr)
	require.NotEmpty(t, addr)
}

func testEnvVarOverride(t *testing.T, override string) {
	port := GetString(ServicePort)
	require.Equal(t, override, port)
}

func testMissingPropertyError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = GetString("redis.not-existing")
}

func testAllRequired(t *testing.T) {
	require.NotEmpty(t, GetString(ServicePort))
	require.NotEmpty(t, GetString(RedisStorePool))
	require.NotEmpty(t, GetString(RedisStoreAddr))
	GetString(RedisStorePassword)
	require.NotEmpty(t, GetString(RedisPubSubPool))
	require.NotEmpty(t, GetString(RedisPubSubAddr))
	GetString(RedisPubSubPassword)
}
