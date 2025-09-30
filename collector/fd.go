package collector

import (
	"host-monitor-agent/models"
	"os"
	"strconv"
	"strings"
)

// FDCollector 文件描述符采集器
type FDCollector struct{}

// Collect 采集文件描述符指标
func (f *FDCollector) Collect() (interface{}, error) {
	metrics := models.FDMetrics{}

	// 读取当前进程已使用的文件描述符数
	// /proc/self/fd 目录下的文件数量
	entries, err := os.ReadDir("/proc/self/fd")
	if err == nil {
		metrics.Allocated = uint64(len(entries))
	}

	// 读取进程级别的文件描述符限制
	// /proc/self/limits 中的 "Max open files"
	limitsData, err := os.ReadFile("/proc/self/limits")
	if err == nil {
		lines := strings.Split(string(limitsData), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Max open files") {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					// Soft limit
					if softLimit, err := strconv.ParseUint(fields[3], 10, 64); err == nil {
						metrics.Maximum = softLimit
						break
					}
				}
			}
		}
	}

	return metrics, nil
}