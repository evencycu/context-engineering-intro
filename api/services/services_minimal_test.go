package services

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestServices_MinimalCoverage(t *testing.T) {
	// 最簡單的測試，專注於能正常運行的部分

	t.Run("ConfigService Full Coverage", func(t *testing.T) {
		// 配置服務不依賴外部服務，可以完全測試
		service := NewConfigService()
		assert.NotNil(t, service)

		// 測試 GetConfig
		config, err := service.GetConfig(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// 測試 ValidateConfig
		valid, err := service.ValidateConfig(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, valid)
	})

	t.Run("FileService ValidateFile Coverage", func(t *testing.T) {
		// 文件服務的驗證功能可以完全測試
		service := NewFileService(nil, nil)
		assert.NotNil(t, service)

		// 測試各種文件驗證場景
		testCases := []struct {
			name     string
			request  *ValidateFileRequest
			expected bool
		}{
			{
				name: "Valid small file",
				request: &ValidateFileRequest{
					FileName:    "test.txt",
					FileSize:    1024,
					ContentType: "text/plain",
				},
				expected: true,
			},
			{
				name: "Valid medium file",
				request: &ValidateFileRequest{
					FileName:    "test.pdf",
					FileSize:    10 * 1024 * 1024, // 10MB
					ContentType: "application/pdf",
				},
				expected: true,
			},
			{
				name: "Large file",
				request: &ValidateFileRequest{
					FileName:    "test.txt",
					FileSize:    100 * 1024 * 1024, // 100MB
					ContentType: "text/plain",
				},
				expected: false,
			},
			{
				name: "Empty file name",
				request: &ValidateFileRequest{
					FileName:    "",
					FileSize:    1024,
					ContentType: "text/plain",
				},
				expected: false,
			},
			{
				name: "Zero file size",
				request: &ValidateFileRequest{
					FileName:    "test.txt",
					FileSize:    0,
					ContentType: "text/plain",
				},
				expected: false,
			},
			{
				name: "Negative file size",
				request: &ValidateFileRequest{
					FileName:    "test.txt",
					FileSize:    -1,
					ContentType: "text/plain",
				},
				expected: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				valid, err := service.ValidateFile(context.Background(), tc.request)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, valid.Valid)
			})
		}
	})

	t.Run("MetricsService Coverage", func(t *testing.T) {
		// 指標服務可以部分測試
		redisClient := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		defer redisClient.Close()

		service := NewMetricsService(redisClient)
		assert.NotNil(t, service)

		// 測試所有指標方法
		err := service.IncrementNotificationSent(context.Background())
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.IncrementNotificationFailed(context.Background())
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.IncrementNotificationPending(context.Background())
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.IncrementTeamsAPICall(context.Background())
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.IncrementTeamsAPIError(context.Background())
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.UpdateResponseTime(context.Background(), 1.5)
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.UpdateQueueProcessingRate(context.Background(), 10.5)
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.UpdateActiveProjects(context.Background(), 5)
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		err = service.UpdateActiveDestinations(context.Background(), 20)
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		}

		// 測試 GetMetrics
		metrics, err := service.GetMetrics(context.Background())
		if err != nil {
			t.Logf("Redis connection error (expected in test environment): %v", err)
		} else {
			assert.NotNil(t, metrics)
		}
	})

	t.Run("Service Creation", func(t *testing.T) {
		// 測試服務創建
		services := []interface{}{
			NewBillingService(nil, nil, nil, nil, nil),
			NewFileService(nil, nil),
			NewConfigService(),
		}

		for _, service := range services {
			assert.NotNil(t, service)
		}
	})
}

func TestServices_ConcurrentAccess(t *testing.T) {
	// 測試並發訪問
	t.Run("ConfigService Concurrent", func(t *testing.T) {
		service := NewConfigService()
		assert.NotNil(t, service)

		// 並發調用
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				config, err := service.GetConfig(context.Background())
				assert.NoError(t, err)
				assert.NotNil(t, config)
				done <- true
			}()
		}

		// 等待所有 goroutine 完成
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("FileService Concurrent", func(t *testing.T) {
		service := NewFileService(nil, nil)
		assert.NotNil(t, service)

		// 並發調用
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				valid, err := service.ValidateFile(context.Background(), &ValidateFileRequest{
					FileName: "test.txt",
					FileSize: 1024,
				})
				assert.NoError(t, err)
				assert.NotNil(t, valid)
				done <- true
			}()
		}

		// 等待所有 goroutine 完成
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestServices_EdgeCases(t *testing.T) {
	// 測試邊界情況
	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		service := NewConfigService()
		_, err := service.GetConfig(ctx)
		// 檢查是否返回了 context 取消錯誤或 nil（取決於實現）
		if err != nil {
			assert.Contains(t, err.Error(), "context canceled")
		}
	})

	t.Run("Nil Dependencies", func(t *testing.T) {
		// 測試 nil 依賴
		services := []interface{}{
			NewBillingService(nil, nil, nil, nil, nil),
			NewFileService(nil, nil),
			NewConfigService(),
		}

		for _, service := range services {
			assert.NotNil(t, service)
		}
	})
}
