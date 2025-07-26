// Package mobile 提供给 Android 客户端使用的接口
package mobile

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/your-username/singbox-xboard-client/internal/config"
	"github.com/your-username/singbox-xboard-client/internal/singbox"
	"github.com/your-username/singbox-xboard-client/internal/subscription"
)

// Client 是给 Android 使用的客户端接口
type Client struct {
	config     *config.Config
	singbox    *singbox.Manager
	subManager *subscription.Manager
	mu         sync.RWMutex
	isRunning  bool
}

// NewClient 创建新的客户端实例
func NewClient() *Client {
	return &Client{
		config:     config.DefaultConfig(),
		subManager: subscription.NewManager(),
	}
}

// Initialize 初始化客户端
func (c *Client) Initialize(configPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 加载配置
	if configPath != "" {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("加载配置失败: %v", err)
		}
		c.config = cfg
	}

	// 初始化 singbox 管理器
	c.singbox = singbox.NewManager(c.config)

	// 初始化订阅管理器
	if err := c.subManager.Initialize(c.config); err != nil {
		return fmt.Errorf("初始化订阅失败: %v", err)
	}

	// 设置订阅更新回调
	c.subManager.OnUpdate(func(singboxConfig map[string]interface{}) {
		c.singbox.UpdateConfig(singboxConfig)
	})

	return nil
}

// Start 启动 VPN 服务
func (c *Client) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isRunning {
		return fmt.Errorf("服务已在运行")
	}

	if err := c.singbox.Start(); err != nil {
		return fmt.Errorf("启动失败: %v", err)
	}

	c.isRunning = true
	return nil
}

// Stop 停止 VPN 服务
func (c *Client) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isRunning {
		return fmt.Errorf("服务未运行")
	}

	if err := c.singbox.Stop(); err != nil {
		return fmt.Errorf("停止失败: %v", err)
	}

	c.isRunning = false
	return nil
}

// UpdateSubscription 更新订阅
func (c *Client) UpdateSubscription(url string) error {
	return c.subManager.UpdateSubscription(url)
}

// RefreshSubscription 刷新当前订阅
func (c *Client) RefreshSubscription() error {
	return c.subManager.RefreshSubscription()
}

// GetStatus 获取当前状态
func (c *Client) GetStatus() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	upload, download, uptime := c.singbox.GetStats()
	
	status := map[string]interface{}{
		"running":  c.isRunning,
		"upload":   upload,
		"download": download,
		"uptime":   uptime.Seconds(),
	}

	// 添加用户信息
	if userInfo := c.subManager.GetUserInfo(); userInfo != nil {
		status["user"] = map[string]interface{}{
			"email":       userInfo.Email,
			"upload":      userInfo.Upload,
			"download":    userInfo.Download,
			"total":       userInfo.Total,
			"expire_time": userInfo.ExpireTime,
		}
	}

	data, _ := json.Marshal(status)
	return string(data)
}

// GetNodes 获取节点列表
func (c *Client) GetNodes() string {
	nodes := c.subManager.GetNodes()
	data, _ := json.Marshal(nodes)
	return string(data)
}

// SelectNode 选择节点
func (c *Client) SelectNode(nodeTag string) error {
	// TODO: 实现节点选择逻辑
	// 这需要修改 singbox 配置中的 selector outbound
	return fmt.Errorf("节点选择功能尚未实现")
}

// GetConfig 获取当前配置
func (c *Client) GetConfig() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, _ := json.Marshal(c.config)
	return string(data)
}

// UpdateConfig 更新配置
func (c *Client) UpdateConfig(configJSON string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var cfg config.Config
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	c.config = &cfg
	
	// 保存配置
	if err := config.Save(c.config, config.GetDefaultConfigPath()); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	// 重新初始化
	return c.Initialize("")
}

// GetLogs 获取日志（最近的N条）
func (c *Client) GetLogs(limit int) string {
	// TODO: 实现日志获取
	return "[]"
}

// TestConnection 测试连接
func (c *Client) TestConnection() string {
	result := map[string]interface{}{
		"success": false,
		"message": "测试中...",
		"delay":   -1,
	}

	// TODO: 实现连接测试
	// 可以通过 HTTP 请求测试代理是否工作

	data, _ := json.Marshal(result)
	return string(data)
}

// Version 获取版本信息
func Version() string {
	return "1.0.0"
}