package collector

import (
	"host-monitor-agent/models"

	"github.com/shirou/gopsutil/v3/load"
)

// LoadCollector 负载指标采集器
type LoadCollector struct{}

// Collect 采集负载指标
func (l *LoadCollector) Collect() (interface{}, error) {
	loadAvg, err := load.Avg()
	if err != nil {
		return models.LoadMetrics{}, err
	}

	return models.LoadMetrics{
		Load1:  loadAvg.Load1,
		Load5:  loadAvg.Load5,
		Load15: loadAvg.Load15,
	}, nil
}