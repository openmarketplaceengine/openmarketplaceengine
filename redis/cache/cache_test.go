package cache

import (
	"context"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/config"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	err := config.Read()
	require.NoError(t, err)

	c := client.NewStoreClient()
	name := "test"
	cache := NewCache(c, name)
	require.Equal(t, cache.GetName(), "cache:test")
	t.Run("testPut", func(t *testing.T) {
		testPut(t, cache)
	})

	t.Run("testRemove", func(t *testing.T) {
		testRemove(t, cache)
	})
}

func testPut(t *testing.T, c Cache) {
	ctx := context.Background()

	key := uuid.New().String()
	value := uuid.New().String()

	containsBefore := c.ContainsKey(ctx, key)
	require.False(t, containsBefore)

	err := c.Put(ctx, key, value)
	require.NoError(t, err)

	containsAfter := c.ContainsKey(ctx, key)
	require.True(t, containsAfter)

	retrieved, err := c.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, retrieved)
}

func testRemove(t *testing.T, c Cache) {
	ctx := context.Background()

	key := uuid.New().String()

	err := c.Put(ctx, key, "anything")
	containsBefore := c.ContainsKey(ctx, key)
	require.True(t, containsBefore)
	require.NoError(t, err)

	err = c.Remove(ctx, key)
	require.NoError(t, err)

	containsAfter := c.ContainsKey(ctx, key)
	require.False(t, containsAfter)
	_, err = c.Get(ctx, key)
	require.Error(t, err)
}
