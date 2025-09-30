package collector

import (
	"host-monitor-agent/models"
	"math"

	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryCollector 内存指标采集器
type MemoryCollector struct{}

// Collect 采集内存指标
func (m *MemoryCollector) Collect() (interface{}, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return models.MemoryMetrics{}, err
	}

	// 转换为 GB (保留1位小数)
	totalGB := math.Round(float64(vmStat.Total)/1024/1024/1024*10) / 10
	usedGB := math.Round(float64(vmStat.Used)/1024/1024/1024*10) / 10

	return models.MemoryMetrics{
		Total:        totalGB,
		Used:         usedGB,
		UsagePercent: math.Round(vmStat.UsedPercent*10) / 10,
	}, nil
}