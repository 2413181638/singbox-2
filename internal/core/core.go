package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"singbox-xboard-client/internal/config"
	"singbox-xboard-client/internal/database"
	"singbox-xboard-client/internal/xboard"
)

type Core struct {
	config    *config.Config
	cmd       *exec.Cmd
	isRunning bool
	mutex     sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

type SingBoxConfig struct {
	Log        LogConfig         `json:"log"`
	DNS        DNSConfig         `json:"dns"`
	Inbounds   []InboundConfig   `json:"inbounds"`
	Outbounds  []OutboundConfig  `json:"outbounds"`
	Route      RouteConfig       `json:"route"`
	Experimental ExperimentalConfig `json:"experimental"`
}

type LogConfig struct {
	Level     string `json:"level"`
	Output    string `json:"output"`
	Timestamp bool   `json:"timestamp"`
}

type DNSConfig struct {
	Servers []DNSServer `json:"servers"`
	Rules   []DNSRule   `json:"rules"`
}

type DNSServer struct {
	Tag     string `json:"tag"`
	Address string `json:"address"`
}

type DNSRule struct {
	Server string   `json:"server"`
	Query  []string `json:"query"`
}

type InboundConfig struct {
	Type   string      `json:"type"`
	Tag    string      `json:"tag"`
	Listen string      `json:"listen"`
	Port   int         `json:"listen_port"`
	Users  []UserConfig `json:"users,omitempty"`
}

type UserConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type OutboundConfig struct {
	Type           string            `json:"type"`
	Tag            string            `json:"tag"`
	Server         string            `json:"server,omitempty"`
	ServerPort     int               `json:"server_port,omitempty"`
	Method         string            `json:"method,omitempty"`
	Password       string            `json:"password,omitempty"`
	UUID           string            `json:"uuid,omitempty"`
	Security       string            `json:"security,omitempty"`
	AlterId        int               `json:"alter_id,omitempty"`
	Network        string            `json:"network,omitempty"`
	TLS            *TLSConfig        `json:"tls,omitempty"`
	Transport      *TransportConfig  `json:"transport,omitempty"`
	Multiplex      *MultiplexConfig  `json:"multiplex,omitempty"`
}

type TLSConfig struct {
	Enabled    bool   `json:"enabled"`
	ServerName string `json:"server_name,omitempty"`
	Insecure   bool   `json:"insecure,omitempty"`
}

type TransportConfig struct {
	Type    string            `json:"type"`
	Path    string            `json:"path,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type MultiplexConfig struct {
	Enabled  bool `json:"enabled"`
	Protocol string `json:"protocol"`
	MaxConnections int `json:"max_connections"`
	MinStreams int `json:"min_streams"`
	MaxStreams int `json:"max_streams"`
}

type RouteConfig struct {
	Rules []RouteRule `json:"rules"`
}

type RouteRule struct {
	Outbound string   `json:"outbound"`
	Domain   []string `json:"domain,omitempty"`
	DomainSuffix []string `json:"domain_suffix,omitempty"`
	IPCidr   []string `json:"ip_cidr,omitempty"`
}

type ExperimentalConfig struct {
	ClashAPI ClashAPIConfig `json:"clash_api"`
}

type ClashAPIConfig struct {
	ExternalController string `json:"external_controller"`
	Secret             string `json:"secret"`
}

func New(config *config.Config) *Core {
	ctx, cancel := context.WithCancel(context.Background())
	return &Core{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Core) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isRunning {
		return fmt.Errorf("core is already running")
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(c.config.SingBox.ConfigPath); os.IsNotExist(err) {
		// 创建默认配置
		if err := c.createDefaultConfig(); err != nil {
			return fmt.Errorf("failed to create default config: %v", err)
		}
	}

	// 启动 SingBox
	c.cmd = exec.CommandContext(c.ctx, "sing-box", "run", "-c", c.config.SingBox.ConfigPath)
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sing-box: %v", err)
	}

	c.isRunning = true
	
	// 监控进程
	go c.monitor()

	return nil
}

func (c *Core) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isRunning {
		return nil
	}

	c.cancel()
	
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill sing-box process: %v", err)
		}
	}

	c.isRunning = false
	return nil
}

func (c *Core) IsRunning() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isRunning
}

func (c *Core) Restart() error {
	if err := c.Stop(); err != nil {
		return err
	}
	
	// 等待进程完全停止
	time.Sleep(2 * time.Second)
	
	return c.Start()
}

func (c *Core) UpdateConfig(servers []database.Server) error {
	config := c.generateConfig(servers)
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(c.config.SingBox.ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	// 如果正在运行，重启以应用新配置
	if c.IsRunning() {
		return c.Restart()
	}

	return nil
}

func (c *Core) createDefaultConfig() error {
	config := &SingBoxConfig{
		Log: LogConfig{
			Level:     c.config.LogLevel,
			Output:    c.config.SingBox.LogPath,
			Timestamp: true,
		},
		DNS: DNSConfig{
			Servers: []DNSServer{
				{Tag: "google", Address: "8.8.8.8"},
				{Tag: "cloudflare", Address: "1.1.1.1"},
			},
		},
		Inbounds: []InboundConfig{
			{
				Type:   "mixed",
				Tag:    "mixed-in",
				Listen: "127.0.0.1",
				Port:   7890,
			},
		},
		Outbounds: []OutboundConfig{
			{
				Type: "direct",
				Tag:  "direct",
			},
			{
				Type: "block",
				Tag:  "block",
			},
		},
		Route: RouteConfig{
			Rules: []RouteRule{
				{
					Outbound: "direct",
					Domain:   []string{"localhost"},
				},
			},
		},
		Experimental: ExperimentalConfig{
			ClashAPI: ClashAPIConfig{
				ExternalController: fmt.Sprintf("127.0.0.1:%d", c.config.SingBox.APIPort),
				Secret:             "",
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.config.SingBox.ConfigPath, data, 0644)
}

func (c *Core) generateConfig(servers []database.Server) *SingBoxConfig {
	config := &SingBoxConfig{
		Log: LogConfig{
			Level:     c.config.LogLevel,
			Output:    c.config.SingBox.LogPath,
			Timestamp: true,
		},
		DNS: DNSConfig{
			Servers: []DNSServer{
				{Tag: "google", Address: "8.8.8.8"},
				{Tag: "cloudflare", Address: "1.1.1.1"},
			},
		},
		Inbounds: []InboundConfig{
			{
				Type:   "mixed",
				Tag:    "mixed-in",
				Listen: "127.0.0.1",
				Port:   7890,
			},
		},
		Outbounds: []OutboundConfig{
			{
				Type: "direct",
				Tag:  "direct",
			},
			{
				Type: "block",
				Tag:  "block",
			},
		},
		Route: RouteConfig{
			Rules: []RouteRule{
				{
					Outbound: "direct",
					Domain:   []string{"localhost"},
				},
			},
		},
		Experimental: ExperimentalConfig{
			ClashAPI: ClashAPIConfig{
				ExternalController: fmt.Sprintf("127.0.0.1:%d", c.config.SingBox.APIPort),
				Secret:             "",
			},
		},
	}

	// 添加服务器配置
	for _, server := range servers {
		if !server.IsActive {
			continue
		}

		outbound := c.convertServerToOutbound(server)
		if outbound != nil {
			config.Outbounds = append(config.Outbounds, *outbound)
		}
	}

	return config
}

func (c *Core) convertServerToOutbound(server database.Server) *OutboundConfig {
	var settings map[string]interface{}
	if err := json.Unmarshal([]byte(server.Settings), &settings); err != nil {
		return nil
	}

	outbound := &OutboundConfig{
		Type:       server.ServerType,
		Tag:        fmt.Sprintf("server-%d", server.ID),
		Server:     server.Host,
		ServerPort: server.Port,
	}

	switch server.ServerType {
	case "shadowsocks":
		if method, ok := settings["method"].(string); ok {
			outbound.Method = method
		}
		if password, ok := settings["password"].(string); ok {
			outbound.Password = password
		}

	case "vmess":
		if uuid, ok := settings["uuid"].(string); ok {
			outbound.UUID = uuid
		}
		if security, ok := settings["security"].(string); ok {
			outbound.Security = security
		}
		if alterId, ok := settings["alter_id"].(float64); ok {
			outbound.AlterId = int(alterId)
		}
		if network, ok := settings["network"].(string); ok {
			outbound.Network = network
		}

	case "vless":
		if uuid, ok := settings["uuid"].(string); ok {
			outbound.UUID = uuid
		}

	case "trojan":
		if password, ok := settings["password"].(string); ok {
			outbound.Password = password
		}
	}

	// TLS 配置
	if tls, ok := settings["tls"].(map[string]interface{}); ok {
		outbound.TLS = &TLSConfig{
			Enabled: true,
		}
		if serverName, ok := tls["server_name"].(string); ok {
			outbound.TLS.ServerName = serverName
		}
		if insecure, ok := tls["insecure"].(bool); ok {
			outbound.TLS.Insecure = insecure
		}
	}

	return outbound
}

func (c *Core) monitor() {
	if c.cmd != nil {
		c.cmd.Wait()
		c.mutex.Lock()
		c.isRunning = false
		c.mutex.Unlock()
	}
}

func (c *Core) GetStats() (map[string]interface{}, error) {
	if !c.IsRunning() {
		return nil, fmt.Errorf("core is not running")
	}

	url := fmt.Sprintf("http://127.0.0.1:%d/traffic", c.config.SingBox.APIPort)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return stats, nil
}