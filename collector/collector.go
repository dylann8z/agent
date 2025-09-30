package collector

import "host-monitor-agent/models"

// Collector 指标采集器接口
type Collector interface {
	Collect() (interface{}, error)
}

// MetricsCollector 所有指标采集器的管理器
type MetricsCollector struct {
	hostInfoCollector Collector
	cpuCollector      Collector
	memoryCollector   Collector
	diskCollector     Collector
	loadCollector     Collector
	tcpCollector      Collector
	fdCollector       Collector
	networkCollector  Collector
	securityCollector Collector
}

// NewMetricsCollector 创建指标采集器管理器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		hostInfoCollector: &HostInfoCollector{},
		cpuCollector:      &CPUCollector{},
		memoryCollector:   &MemoryCollector{},
		diskCollector:     &DiskCollector{},
		loadCollector:     &LoadCollector{},
		tcpCollector:      &TCPCollector{},
		fdCollector:       &FDCollector{},
		networkCollector:  &NetworkCollector{},
		securityCollector: &SecurityCollector{},
	}
}

// CollectAll 采集所有指标
func (mc *MetricsCollector) CollectAll() (*models.HostMetrics, error) {
	metrics := &models.HostMetrics{}

	// 采集主机信息
	if hostInfo, err := mc.hostInfoCollector.Collect(); err == nil {
		info := hostInfo.(HostInfo)
		metrics.Hostname = info.Hostname
		metrics.IntranetIPs = info.IntranetIPs
		metrics.OS = info.OS
		metrics.KernelVersion = info.KernelVersion
		metrics.Timezone = info.Timezone
		metrics.Uptime = info.Uptime
	}

	// 采集CPU
	if cpu, err := mc.cpuCollector.Collect(); err == nil {
		metrics.CPU = cpu.(models.CPUMetrics)
	}

	// 采集内存
	if mem, err := mc.memoryCollector.Collect(); err == nil {
		metrics.Memory = mem.(models.MemoryMetrics)
	}

	// 采集磁盘
	if disk, err := mc.diskCollector.Collect(); err == nil {
		metrics.Disk = disk.([]models.DiskMetrics)
	}

	// 采集负载
	if load, err := mc.loadCollector.Collect(); err == nil {
		metrics.Load = load.(models.LoadMetrics)
	}

	// 采集TCP连接
	if tcp, err := mc.tcpCollector.Collect(); err == nil {
		metrics.TCP = tcp.(models.TCPMetrics)
	}

	// 采集文件描述符
	if fd, err := mc.fdCollector.Collect(); err == nil {
		metrics.FileDescriptor = fd.(models.FDMetrics)
	}

	// 采集网络流量
	if network, err := mc.networkCollector.Collect(); err == nil {
		metrics.Network = network.([]models.NetworkMetrics)
	}

	// 采集安全指标
	if security, err := mc.securityCollector.Collect(); err == nil {
		metrics.Security = security.(models.SecurityMetrics)
	}

	return metrics, nil
}