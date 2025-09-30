package api

import (
	"host-monitor-agent/collector"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	metricsCollector *collector.MetricsCollector
}

// NewHandler 创建API处理器
func NewHandler() *Handler {
	return &Handler{
		metricsCollector: collector.NewMetricsCollector(),
	}
}

// GetMetrics 获取监控指标
func (h *Handler) GetMetrics(c *gin.Context) {
	metrics, err := h.metricsCollector.CollectAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 设置时间戳 (格式: 2025-09-30 07:55:16 UTC)
	metrics.Timestamp = time.Now().UTC().Format("2006-01-02 15:04:05 MST")

	c.JSON(http.StatusOK, metrics)
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now(),
	})
}