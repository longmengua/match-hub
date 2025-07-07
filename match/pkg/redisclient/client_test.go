package redisclient_test

import (
	"context"
	"log"
	"testing"
	"time"

	"match/pkg/redisclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisClient_Single(t *testing.T) {
	ctx := context.Background()

	client, err := redisclient.New("single", []string{"redis-signle.docker-compose-gui.orb.local:6379"}, "", 0)
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
		"redis-1.docker-compose-gui.orb.local:6379",
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

func TestRedisClient_GetWithFallback(t *testing.T) {
	ctx := context.Background()

	client, err := redisclient.New("single", []string{"redis-signle.redis.orb.local:6379"}, "", 0)
	require.NoError(t, err)

	key := "test:key:fallback"
	val := "from_db"
	exp := time.Minute

	// 確保起始 Redis 無資料
	_ = client.Del(ctx, key)

	// counter 模擬 DB 查詢次數
	callCount := 0

	fallback := func() (string, error) {
		callCount++
		time.Sleep(500 * time.Millisecond) // 模擬耗時操作
		log.Printf("DB Get: %s", val)
		return val, nil
	}

	// 啟動多個 goroutine 同時呼叫
	n := 100
	results := make(chan string, n)

	for i := 0; i < n; i++ {
		go func() {
			v, err := client.GetWithFallback(ctx, key, exp, fallback)
			require.NoError(t, err)
			results <- v
		}()
	}

	// 收集結果
	for i := 0; i < n; i++ {
		got := <-results
		assert.Equal(t, val, got)
	}

	// 確保 fallback 只被呼叫一次（singleflight 有效）
	assert.Equal(t, 1, callCount)

	// 確保之後從 Redis 快取中讀到資料
	cached, err := client.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, val, cached)
}
