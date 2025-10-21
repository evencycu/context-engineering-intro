package integration

import (
	"context"
	"testing"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/api/services"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_MetricsService_WithRedis(t *testing.T) {
	// 集成測試：指標服務與 Redis
	// 注意：此測試需要 Redis 實例運行

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	// 測試 Redis 連接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	// 創建指標服務
	metricsService := services.NewMetricsService(redisClient)
	assert.NotNil(t, metricsService)

	// 測試基本功能
	err = metricsService.IncrementNotificationSent(ctx)
	assert.NoError(t, err)

	err = metricsService.IncrementNotificationFailed(ctx)
	assert.NoError(t, err)

	err = metricsService.IncrementNotificationPending(ctx)
	assert.NoError(t, err)

	// 測試更新方法
	err = metricsService.UpdateResponseTime(ctx, 1.5)
	assert.NoError(t, err)

	err = metricsService.UpdateQueueProcessingRate(ctx, 10.5)
	assert.NoError(t, err)

	err = metricsService.UpdateActiveProjects(ctx, 5)
	assert.NoError(t, err)

	err = metricsService.UpdateActiveDestinations(ctx, 20)
	assert.NoError(t, err)

	// 測試獲取指標
	metrics, err := metricsService.GetMetrics(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.System)
	assert.NotNil(t, metrics.Redis)
	assert.NotNil(t, metrics.Business)
}

func TestIntegration_MetricsService_ConcurrentAccess(t *testing.T) {
	// 集成測試：並發訪問
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	metricsService := services.NewMetricsService(redisClient)
	assert.NotNil(t, metricsService)

	// 並發調用
	done := make(chan bool, 20)
	for i := 0; i < 20; i++ {
		go func() {
			err := metricsService.IncrementNotificationSent(ctx)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestIntegration_MetricsService_Performance(t *testing.T) {
	// 集成測試：性能測試
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	metricsService := services.NewMetricsService(redisClient)
	assert.NotNil(t, metricsService)

	start := time.Now()

	// 執行大量操作
	for i := 0; i < 1000; i++ {
		err := metricsService.IncrementNotificationSent(ctx)
		assert.NoError(t, err)
	}

	duration := time.Since(start)
	t.Logf("Executed 1000 operations in %v", duration)
	assert.Less(t, duration, 5*time.Second) // 應該在5秒內完成
}

func TestIntegration_MetricsService_DataPersistence(t *testing.T) {
	// 集成測試：數據持久化
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	metricsService := services.NewMetricsService(redisClient)
	assert.NotNil(t, metricsService)

	// 設置一些數據
	err = metricsService.UpdateActiveProjects(ctx, 10)
	assert.NoError(t, err)

	err = metricsService.UpdateActiveDestinations(ctx, 50)
	assert.NoError(t, err)

	// 獲取指標並驗證數據
	metrics, err := metricsService.GetMetrics(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.Business)
}

func TestIntegration_MetricsService_ErrorRecovery(t *testing.T) {
	// 集成測試：錯誤恢復
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	metricsService := services.NewMetricsService(redisClient)
	assert.NotNil(t, metricsService)

	// 測試正常操作
	err = metricsService.IncrementNotificationSent(ctx)
	assert.NoError(t, err)

	// 測試錯誤處理
	invalidCtx, invalidCancel := context.WithCancel(context.Background())
	invalidCancel() // 取消上下文

	err = metricsService.IncrementNotificationSent(invalidCtx)
	assert.Error(t, err) // 應該返回錯誤
}

func TestIntegration_MetricsService_LoadTesting(t *testing.T) {
	// 集成測試：負載測試
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	metricsService := services.NewMetricsService(redisClient)
	assert.NotNil(t, metricsService)

	// 負載測試
	start := time.Now()

	// 並發執行多個操作
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				err := metricsService.IncrementNotificationSent(ctx)
				if err != nil {
					t.Logf("Error in load test: %v", err)
					break
				}
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 100; i++ {
		<-done
	}

	duration := time.Since(start)
	t.Logf("Load test completed in %v", duration)
	assert.Less(t, duration, 30*time.Second) // 應該在30秒內完成
}
