package client

import (
	"context"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/config"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {

	err := config.Read()
	require.NoError(t, err)

	t.Run("testCreateClients", func(t *testing.T) {
		testCreateClients(t)
	})
	t.Run("testCommand", func(t *testing.T) {
		testCommand(t)
	})
	t.Run("testZ", func(t *testing.T) {
		testZ(t)
	})
}

func testCreateClients(t *testing.T) {
	client := NewStoreClient()
	require.NotNil(t, client)
	clientPool := NewStoreClientPool()
	require.NotNil(t, clientPool)

	pubSubClient := NewPubSubClient()
	require.NotNil(t, pubSubClient)
	pubSubClientPool := NewPubSubClientPool()
	require.NotNil(t, pubSubClientPool)
}

func testCommand(t *testing.T) {
	client := NewStoreClientPool()
	require.NotNil(t, client)

	ctx := context.Background()

	val, err := client.Set(ctx, "foo", "someval", 0).Result()
	require.NoError(t, err)
	require.Equal(t, "OK", val)

	val, err = client.Get(ctx, "foo").Result()
	require.NoError(t, err)
	require.Equal(t, "someval", val)
}

func testZ(t *testing.T) {
	client := NewStoreClientPool()
	require.NotNil(t, client)

	ctx := context.Background()

	key := "testZ"
	client.Unlink(ctx, key)

	z := redis.Z{Score: 1, Member: "some-value"}
	val, err := client.ZAddArgs(ctx, key, redis.ZAddArgs{Ch: true, Members: []redis.Z{z}}).Result()
	require.NoError(t, err)
	expected := int64(1)
	if val != expected {
		assert.Fail(t, "expected to add value")
	}
}
