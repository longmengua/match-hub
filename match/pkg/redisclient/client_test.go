package redisclient_test

import (
	"context"
	"testing"
	"time"

	"match/pkg/redisclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 請先確保你在本地或測試環境中有一個 Redis single instance 以及 Redis cluster 可用

func TestRedisClient_Single(t *testing.T) {
	ctx := context.Background()

	client, err := redisclient.New("single", []string{"redis.docker-compose-gui.orb.local:6390"}, "", 0)
	require.NoError(t, err)

	key := "test:key:single"
	value := "hello_single"

	// Set
	err = client.Set(ctx, key, value, time.Minute)
	require.NoError(t, err)

	// Get
	got, err := client.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value, got)

	// Del
	err = client.Del(ctx, key)
	require.NoError(t, err)
}

func TestRedisClient_Cluster(t *testing.T) {
	ctx := context.Background()

	// 請確認這些地址是你的 cluster 節點
	clusterAddrs := []string{
		"redis-master.docker-compose-gui.orb.local:6379",
	}
	client, err := redisclient.New("cluster", clusterAddrs, "", 0)
	require.NoError(t, err)

	key := "test:key:cluster"
	value := "hello_cluster"

	// Set
	err = client.Set(ctx, key, value, time.Minute)
	require.NoError(t, err)

	// Get
	got, err := client.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value, got)

	// Del
	err = client.Del(ctx, key)
	require.NoError(t, err)
}
