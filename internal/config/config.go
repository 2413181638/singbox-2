package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Inbound  InboundConfig  `yaml:"inbound"`
	Outbound OutboundConfig `yaml:"outbound"`
	DNS      DNSConfig      `yaml:"dns"`
	Log      LogConfig      `yaml:"log"`
}

type InboundConfig struct {
	Type string `yaml:"type"`
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type OutboundConfig struct {
	Type     string            `yaml:"type"`
	Server   string            `yaml:"server"`
	Port     int               `yaml:"port"`
	Username string            `yaml:"username,omitempty"`
	Password string            `yaml:"password,omitempty"`
	Method   string            `yaml:"method,omitempty"`
	Options  map[string]string `yaml:"options,omitempty"`
}

type DNSConfig struct {
	Servers []string `yaml:"servers"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file,omitempty"`
}

func Load(filename string) (*Config, error) {
	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		defaultConfig := getDefaultConfig()
		if err := Save(filename, defaultConfig); err != nil {
			return nil, fmt.Errorf("创建默认配置失败: %w", err)
		}
		return defaultConfig, nil
	}
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	
	return &config, nil
}

func Save(filename string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	
	return nil
}

func getDefaultConfig() *Config {
	return &Config{
		Inbound: InboundConfig{
			Type: "socks",
			Port: 1080,
			Host: "127.0.0.1",
		},
		Outbound: OutboundConfig{
			Type:   "direct",
			Server: "",
			Port:   0,
		},
		DNS: DNSConfig{
			Servers: []string{"8.8.8.8", "1.1.1.1"},
		},
		Log: LogConfig{
			Level: "info",
		},
	}
}