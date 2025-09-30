package config

import "time"

// Config 配置
type Config struct {
	Server    ServerConfig    `json:"server"`
	Collector CollectorConfig `json:"collector"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// CollectorConfig 采集器配置
type CollectorConfig struct {
	Interval time.Duration `json:"interval"` // 采集间隔（秒）
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Collector: CollectorConfig{
			Interval: 10 * time.Second, // 默认10秒采集一次
		},
	}
}