package collector

import (
	"bufio"
	"host-monitor-agent/models"
	"os"
	"strings"
)

// SecurityCollector 安全指标采集器
type SecurityCollector struct{}

// Collect 采集安全指标
func (s *SecurityCollector) Collect() (interface{}, error) {
	metrics := models.SecurityMetrics{}

	// 统计登录失败次数
	loginFailures, err := countLoginFailures()
	if err == nil {
		metrics.LoginFailures = loginFailures
	}

	return metrics, nil
}

// countLoginFailures 统计登录失败次数
func countLoginFailures() (uint64, error) {
	// 尝试读取不同的日志文件
	logFiles := []string{
		"/var/log/auth.log",      // Debian/Ubuntu
		"/var/log/secure",         // RHEL/CentOS/Amazon Linux
		"/var/log/faillog",        // 通用
	}

	var count uint64

	for _, logFile := range logFiles {
		file, err := os.Open(logFile)
		if err != nil {
			continue // 文件不存在或无权限，尝试下一个
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			// 匹配常见的登录失败关键字
			if strings.Contains(line, "Failed password") ||
				strings.Contains(line, "authentication failure") ||
				strings.Contains(line, "Invalid user") ||
				strings.Contains(line, "Failed login") {
				count++
			}
		}

		// 如果成功读取了一个文件，就返回结果
		if scanner.Err() == nil {
			return count, nil
		}
	}

	// 如果所有文件都无法读取，返回0
	return 0, nil
}