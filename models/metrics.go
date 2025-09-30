package models

// HostMetrics 主机所有监控指标
type HostMetrics struct {
	Timestamp      string           `json:"timestamp"`
	Hostname       string           `json:"hostname"`
	IntranetIPs    []string         `json:"intranet_ips"`
	OS             string           `json:"os"`              // 操作系统发行版
	KernelVersion  string           `json:"kernel_version"`  // 内核版本
	Timezone       string           `json:"timezone"`        // 时区
	Uptime         string           `json:"uptime"`          // 运行时间
	CPU            CPUMetrics       `json:"cpu"`
	Memory         MemoryMetrics    `json:"memory"`
	Disk           []DiskMetrics    `json:"disk"`
	Load           LoadMetrics      `json:"load"`
	TCP            TCPMetrics       `json:"tcp"`
	FileDescriptor FDMetrics        `json:"file_descriptor"`
	Network        []NetworkMetrics `json:"network"`
	Security       SecurityMetrics  `json:"security"`
}

// CPUMetrics CPU监控指标
type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"`
	CoreCount    int     `json:"core_count"`
}

// MemoryMetrics 内存监控指标
type MemoryMetrics struct {
	Total        float64 `json:"total_gb"`
	Used         float64 `json:"used_gb"`
	UsagePercent float64 `json:"usage_percent"`
}

// DiskMetrics 磁盘监控指标
type DiskMetrics struct {
	MountPoint   string  `json:"mount_point"`
	Total        float64 `json:"total_gb"`
	Used         float64 `json:"used_gb"`
	UsagePercent float64 `json:"usage_percent"`
}

// LoadMetrics 负载监控指标
type LoadMetrics struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// IOMetrics IO监控指标（已移除，IO统计意义不大）

// TCPMetrics TCP连接监控指标
type TCPMetrics struct {
	Established uint64 `json:"established"`
	SynSent     uint64 `json:"syn_sent"`
	SynRecv     uint64 `json:"syn_recv"`
	FinWait1    uint64 `json:"fin_wait1"`
	FinWait2    uint64 `json:"fin_wait2"`
	TimeWait    uint64 `json:"time_wait"`
	Close       uint64 `json:"close"`
	CloseWait   uint64 `json:"close_wait"`
	LastAck     uint64 `json:"last_ack"`
	Listen      uint64 `json:"listen"`
	Closing     uint64 `json:"closing"`
	Total       uint64 `json:"total"`
}

// FDMetrics 文件描述符监控指标
type FDMetrics struct {
	Allocated uint64 `json:"allocated"` // 已分配的文件描述符数
	Maximum   uint64 `json:"maximum"`   // 最大文件描述符数
}

// NetworkMetrics 网络流量监控指标
type NetworkMetrics struct {
	Interface string `json:"interface"`
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
}

// SecurityMetrics 安全监控指标
type SecurityMetrics struct {
	LoginFailures uint64 `json:"login_failures"` // 登录失败次数
}