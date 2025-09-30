package cache

import (
	"host-monitor-agent/collector"
	"host-monitor-agent/models"
	"log"
	"sync"
	"time"
)

// MetricsCache 监控指标缓存
type MetricsCache struct {
	metrics   *models.HostMetrics
	mutex     sync.RWMutex
	collector *collector.MetricsCollector
	interval  time.Duration
	stopChan  chan struct{}
}

// NewMetricsCache 创建缓存实例
func NewMetricsCache(interval time.Duration) *MetricsCache {
	return &MetricsCache{
		collector: collector.NewMetricsCollector(),
		interval:  interval,
		stopChan:  make(chan struct{}),
	}
}

// Start 启动后台定时采集
func (c *MetricsCache) Start() {
	// 启动时立即采集一次
	c.update()

	// 启动定时器
	ticker := time.NewTicker(c.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.update()
			case <-c.stopChan:
				ticker.Stop()
				log.Println("Metrics cache collector stopped")
				return
			}
		}
	}()

	log.Printf("Metrics cache collector started (interval: %v)", c.interval)
}

// Stop 停止后台采集
func (c *MetricsCache) Stop() {
	close(c.stopChan)
}

// update 更新缓存数据
func (c *MetricsCache) update() {
	metrics, err := c.collector.CollectAll()
	if err != nil {
		log.Printf("Failed to collect metrics: %v", err)
		return
	}

	// 设置时间戳
	metrics.Timestamp = time.Now().UTC().Format("2006-01-02 15:04:05 MST")

	// 更新缓存
	c.mutex.Lock()
	c.metrics = metrics
	c.mutex.Unlock()
}

// Get 获取缓存的监控指标
func (c *MetricsCache) Get() *models.HostMetrics {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.metrics
}