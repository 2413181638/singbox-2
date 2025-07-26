package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	// 基础配置
	LogLevel string `json:"log_level" yaml:"log_level"`
	
	// 订阅配置
	Subscription SubscriptionConfig `json:"subscription" yaml:"subscription"`
	
	// Sing-box 配置
	Singbox SingboxConfig `json:"singbox" yaml:"singbox"`
	
	// UI 配置
	UI UIConfig `json:"ui" yaml:"ui"`
	
	// 规则配置
	Rules RulesConfig `json:"rules" yaml:"rules"`
}

// SubscriptionConfig 订阅配置
type SubscriptionConfig struct {
	URL            string `json:"url" yaml:"url"`                         // 订阅地址
	Token          string `json:"token" yaml:"token"`                     // 认证令牌
	UpdateInterval int    `json:"update_interval" yaml:"update_interval"` // 更新间隔（分钟）
	AutoUpdate     bool   `json:"auto_update" yaml:"auto_update"`         // 自动更新
}

// SingboxConfig sing-box 相关配置
type SingboxConfig struct {
	Inbounds  []InboundConfig  `json:"inbounds" yaml:"inbounds"`   // 入站配置
	Outbounds []OutboundConfig `json:"outbounds" yaml:"outbounds"` // 出站配置
	DNS       DNSConfig        `json:"dns" yaml:"dns"`             // DNS 配置
	Route     RouteConfig      `json:"route" yaml:"route"`         // 路由配置
}

// InboundConfig 入站配置
type InboundConfig struct {
	Type       string `json:"type" yaml:"type"`               // 类型：mixed, tun
	Tag        string `json:"tag" yaml:"tag"`                 // 标签
	Listen     string `json:"listen" yaml:"listen"`           // 监听地址
	ListenPort int    `json:"listen_port" yaml:"listen_port"` // 监听端口
}

// OutboundConfig 出站配置
type OutboundConfig struct {
	Type string `json:"type" yaml:"type"` // 类型：direct, block, selector, urltest
	Tag  string `json:"tag" yaml:"tag"`   // 标签
}

// DNSConfig DNS 配置
type DNSConfig struct {
	Servers []DNSServer `json:"servers" yaml:"servers"` // DNS 服务器列表
	Rules   []DNSRule   `json:"rules" yaml:"rules"`     // DNS 规则
}

// DNSServer DNS 服务器配置
type DNSServer struct {
	Tag     string `json:"tag" yaml:"tag"`         // 标签
	Address string `json:"address" yaml:"address"` // 地址
	Detour  string `json:"detour" yaml:"detour"`   // 出站标签
}

// DNSRule DNS 规则
type DNSRule struct {
	Domain []string `json:"domain,omitempty" yaml:"domain,omitempty"` // 域名列表
	Server string   `json:"server" yaml:"server"`                     // DNS 服务器标签
}

// RouteConfig 路由配置
type RouteConfig struct {
	Rules    []RouteRule `json:"rules" yaml:"rules"`       // 路由规则
	RuleSet  []RuleSet   `json:"rule_set" yaml:"rule_set"` // 规则集
	GeoIP    GeoIPConfig `json:"geoip" yaml:"geoip"`       // GeoIP 配置
	GeoSite  GeoSiteConfig `json:"geosite" yaml:"geosite"`   // GeoSite 配置
}

// RouteRule 路由规则
type RouteRule struct {
	Type     string   `json:"type,omitempty" yaml:"type,omitempty"`         // 规则类型
	Domain   []string `json:"domain,omitempty" yaml:"domain,omitempty"`     // 域名
	IP       []string `json:"ip,omitempty" yaml:"ip,omitempty"`             // IP
	Port     []int    `json:"port,omitempty" yaml:"port,omitempty"`         // 端口
	Protocol []string `json:"protocol,omitempty" yaml:"protocol,omitempty"` // 协议
	RuleSet  []string `json:"rule_set,omitempty" yaml:"rule_set,omitempty"` // 规则集
	Outbound string   `json:"outbound" yaml:"outbound"`                     // 出站标签
}

// RuleSet 规则集配置
type RuleSet struct {
	Tag    string `json:"tag" yaml:"tag"`       // 标签
	Type   string `json:"type" yaml:"type"`     // 类型：local, remote
	Format string `json:"format" yaml:"format"` // 格式：binary, source
	Path   string `json:"path,omitempty" yaml:"path,omitempty"` // 本地路径
	URL    string `json:"url,omitempty" yaml:"url,omitempty"`   // 远程地址
}

// GeoIPConfig GeoIP 配置
type GeoIPConfig struct {
	Path           string `json:"path" yaml:"path"`                       // 数据库路径
	DownloadURL    string `json:"download_url" yaml:"download_url"`       // 下载地址
	DownloadDetour string `json:"download_detour" yaml:"download_detour"` // 下载使用的出站
}

// GeoSiteConfig GeoSite 配置
type GeoSiteConfig struct {
	Path           string `json:"path" yaml:"path"`                       // 数据库路径
	DownloadURL    string `json:"download_url" yaml:"download_url"`       // 下载地址
	DownloadDetour string `json:"download_detour" yaml:"download_detour"` // 下载使用的出站
}

// UIConfig UI 配置
type UIConfig struct {
	Listen string `json:"listen" yaml:"listen"` // 监听地址
	Port   int    `json:"port" yaml:"port"`     // 监听端口
	Secret string `json:"secret" yaml:"secret"` // API 密钥
}

// RulesConfig 规则配置
type RulesConfig struct {
	Mode string `json:"mode" yaml:"mode"` // 规则模式：rule, global, direct
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		LogLevel: "info",
		Subscription: SubscriptionConfig{
			UpdateInterval: 60,
			AutoUpdate:     true,
		},
		Singbox: SingboxConfig{
			Inbounds: []InboundConfig{
				{
					Type:       "mixed",
					Tag:        "mixed-in",
					Listen:     "127.0.0.1",
					ListenPort: 7890,
				},
				{
					Type: "tun",
					Tag:  "tun-in",
				},
			},
			Outbounds: []OutboundConfig{
				{Type: "direct", Tag: "direct"},
				{Type: "block", Tag: "block"},
				{Type: "selector", Tag: "proxy"},
			},
			DNS: DNSConfig{
				Servers: []DNSServer{
					{
						Tag:     "remote",
						Address: "https://1.1.1.1/dns-query",
						Detour:  "proxy",
					},
					{
						Tag:     "local",
						Address: "https://223.5.5.5/dns-query",
						Detour:  "direct",
					},
				},
			},
		},
		UI: UIConfig{
			Listen: "127.0.0.1",
			Port:   9090,
		},
		Rules: RulesConfig{
			Mode: "rule",
		},
	}
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	// 如果没有指定路径，使用默认路径
	if path == "" {
		path = GetDefaultConfigPath()
	}

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := Save(cfg, path); err != nil {
			return nil, fmt.Errorf("创建默认配置失败: %w", err)
		}
		return cfg, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	cfg := &Config{}
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		err = json.Unmarshal(data, cfg)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, cfg)
	default:
		// 尝试 JSON
		if err = json.Unmarshal(data, cfg); err != nil {
			// 尝试 YAML
			err = yaml.Unmarshal(data, cfg)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return cfg, nil
}

// Save 保存配置文件
func Save(cfg *Config, path string) error {
	// 创建目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 序列化配置
	var data []byte
	var err error
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		data, err = json.MarshalIndent(cfg, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(cfg)
	default:
		// 默认使用 JSON
		data, err = json.MarshalIndent(cfg, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetDefaultConfigPath 获取默认配置文件路径
func GetDefaultConfigPath() string {
	// 获取用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		// 使用当前目录
		configDir = "."
	}

	return filepath.Join(configDir, "singbox-xboard", "config.json")
}

// GetConfigDir 获取配置目录
func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "singbox-xboard")
}