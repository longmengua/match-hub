package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"
)

var srv *http.Server

func Start() error {
	router := SetupRouter()

	srv = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("HTTP server listening on :8080")
	return srv.ListenAndServe()
}

func Stop() {
	log.Println("Shutting down HTTP server...")
	// 會等最多 5 秒讓現有連線處理完，超過 5 秒就會強制中止等待（ctx.Done() 觸發），server 被迫終止
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 優雅關閉
	cleanup()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
}

// 優雅關閉 HTTP server
func cleanup() {
	log.Println("Start cleanup tasks...")
	// 這裡可以放你想做的資源關閉邏輯，例如
	// - 關閉 DB 連線
	// - 關閉外部服務的連線
	// - 儲存緩存資料
	// - 停止背景工作（worker）
	// - 紀錄日志等

	time.Sleep(2 * time.Second) // 模擬一些清理動作
	log.Println("Cleanup finished.")
}
