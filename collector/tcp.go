package collector

import (
	"host-monitor-agent/models"

	"github.com/shirou/gopsutil/v3/net"
)

// TCPCollector TCP连接采集器
type TCPCollector struct{}

// Collect 采集TCP连接指标
func (t *TCPCollector) Collect() (interface{}, error) {
	connections, err := net.Connections("tcp")
	if err != nil {
		return models.TCPMetrics{}, err
	}

	metrics := models.TCPMetrics{}

	// 统计各状态的连接数
	for _, conn := range connections {
		metrics.Total++

		switch conn.Status {
		case "ESTABLISHED":
			metrics.Established++
		case "SYN_SENT":
			metrics.SynSent++
		case "SYN_RECV":
			metrics.SynRecv++
		case "FIN_WAIT1":
			metrics.FinWait1++
		case "FIN_WAIT2":
			metrics.FinWait2++
		case "TIME_WAIT":
			metrics.TimeWait++
		case "CLOSE":
			metrics.Close++
		case "CLOSE_WAIT":
			metrics.CloseWait++
		case "LAST_ACK":
			metrics.LastAck++
		case "LISTEN":
			metrics.Listen++
		case "CLOSING":
			metrics.Closing++
		}
	}

	return metrics, nil
}