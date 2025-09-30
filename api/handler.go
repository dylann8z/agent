package api

import (
	"host-monitor-agent/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	metricsCache *cache.MetricsCache
}

// NewHandler 创建API处理器
func NewHandler(metricsCache *cache.MetricsCache) *Handler {
	return &Handler{
		metricsCache: metricsCache,
	}
}

// GetMetrics 获取监控指标（从缓存读取）
func (h *Handler) GetMetrics(c *gin.Context) {
	metrics := h.metricsCache.Get()

	if metrics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "metrics not ready",
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now(),
	})
}