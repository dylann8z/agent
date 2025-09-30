package collector

import (
	"host-monitor-agent/models"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUCollector CPU指标采集器
type CPUCollector struct{}

// Collect 采集CPU指标
func (c *CPUCollector) Collect() (interface{}, error) {
	// 获取CPU使用率（1秒采样）
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return models.CPUMetrics{}, err
	}

	usage := 0.0
	if len(percentages) > 0 {
		usage = percentages[0]
	}

	return models.CPUMetrics{
		UsagePercent: usage,
		CoreCount:    runtime.NumCPU(),
	}, nil
}