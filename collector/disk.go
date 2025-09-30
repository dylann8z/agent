package collector

import (
	"host-monitor-agent/models"
	"math"

	"github.com/shirou/gopsutil/v3/disk"
)

// DiskCollector 磁盘指标采集器
type DiskCollector struct{}

// Collect 采集磁盘指标
func (d *DiskCollector) Collect() (interface{}, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return []models.DiskMetrics{}, err
	}

	var diskMetrics []models.DiskMetrics
	for _, partition := range partitions {
		// 过滤掉不需要监控的分区
		if shouldSkipPartition(partition.Mountpoint) {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		// 只监控容量大于1GB的分区
		if usage.Total < 1024*1024*1024 {
			continue
		}

		// 转换为 GB (保留1位小数)
		totalGB := math.Round(float64(usage.Total)/1024/1024/1024*10) / 10
		usedGB := math.Round(float64(usage.Used)/1024/1024/1024*10) / 10

		diskMetrics = append(diskMetrics, models.DiskMetrics{
			MountPoint:   partition.Mountpoint,
			Total:        totalGB,
			Used:         usedGB,
			UsagePercent: math.Round(usage.UsedPercent*10) / 10,
		})
	}

	return diskMetrics, nil
}

// shouldSkipPartition 判断是否跳过该分区
func shouldSkipPartition(mountPoint string) bool {
	// 跳过的挂载点列表
	skipList := []string{
		"/boot",
		"/boot/efi",
		"/sys",
		"/proc",
		"/dev",
		"/run",
		"/snap",
	}

	for _, skip := range skipList {
		if mountPoint == skip || len(mountPoint) > len(skip) && mountPoint[:len(skip)+1] == skip+"/" {
			return true
		}
	}

	return false
}