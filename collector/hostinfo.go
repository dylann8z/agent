package collector

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

// HostInfoCollector 主机信息采集器
type HostInfoCollector struct{}

// HostInfo 主机信息
type HostInfo struct {
	Hostname      string
	IntranetIPs   []string
	OS            string
	KernelVersion string
	Timezone      string
	Uptime        string
}

// Collect 采集主机信息
func (h *HostInfoCollector) Collect() (interface{}, error) {
	info := HostInfo{
		IntranetIPs: []string{},
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	// 获取系统信息
	hostInfo, err := host.Info()
	if err == nil {
		// 操作系统发行版: "ubuntu 22.04", "centos 7.9"
		info.OS = hostInfo.Platform + " " + hostInfo.PlatformVersion
		// 内核版本
		info.KernelVersion = hostInfo.KernelVersion
		// 运行时间（格式化为人类可读）
		info.Uptime = formatUptime(hostInfo.Uptime)
	}

	// 获取时区
	_, offset := time.Now().Zone()
	zone, _ := time.Now().Zone()
	if zone == "" {
		// 如果获取不到时区名称，使用 UTC+offset 格式
		hours := offset / 3600
		if hours >= 0 {
			info.Timezone = fmt.Sprintf("UTC+%d", hours)
		} else {
			info.Timezone = fmt.Sprintf("UTC%d", hours)
		}
	} else {
		info.Timezone = zone
	}

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return info, err
	}

	// 遍历所有网络接口获取内网IP
	for _, iface := range interfaces {
		// 跳过 loopback 和未启用的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 获取所有IPv4地址（包括内网和公网）
			if ip != nil && ip.To4() != nil {
				// 排除本地回环地址
				if !ip.IsLoopback() {
					info.IntranetIPs = append(info.IntranetIPs, ip.String())
				}
			}
		}
	}

	return info, nil
}

// formatUptime 将秒数格式化为人类可读的运行时间
func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

// isPrivateIP 判断是否为内网IP
func isPrivateIP(ip net.IP) bool {
	// 10.0.0.0/8
	if ip[0] == 10 {
		return true
	}
	// 172.16.0.0/12
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}
	// 192.168.0.0/16
	if ip[0] == 192 && ip[1] == 168 {
		return true
	}
	return false
}