package collector

import (
	"host-monitor-agent/models"

	"github.com/shirou/gopsutil/v3/net"
)

// NetworkCollector 网络流量采集器
type NetworkCollector struct{}

// Collect 采集网络流量指标
func (n *NetworkCollector) Collect() (interface{}, error) {
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return []models.NetworkMetrics{}, err
	}

	var networkMetrics []models.NetworkMetrics
	for _, counter := range ioCounters {
		// 跳过 loopback 接口
		if counter.Name == "lo" {
			continue
		}

		networkMetrics = append(networkMetrics, models.NetworkMetrics{
			Interface: counter.Name,
			BytesSent: counter.BytesSent,
			BytesRecv: counter.BytesRecv,
		})
	}

	return networkMetrics, nil
}