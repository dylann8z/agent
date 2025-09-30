package api

import (
	"host-monitor-agent/cache"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(metricsCache *cache.MetricsCache) *gin.Engine {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	handler := NewHandler(metricsCache)

	// 健康检查
	router.GET("/health", handler.HealthCheck)

	// 获取监控指标
	router.GET("/metrics", handler.GetMetrics)

	return router
}