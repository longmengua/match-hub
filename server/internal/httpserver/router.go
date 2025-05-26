package httpserver

import (
	handler "server/internal/hanlder"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 公用中介層 (middleware)
	// router.Use(SomeMiddleware())

	// 註冊健康檢查路由群組
	healthGroup := router.Group("/health")
	handler.RegisterHealthRoutes(healthGroup)

	// 可以繼續新增更多路由群組...

	return router
}
