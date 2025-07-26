package subscription

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/your-username/singbox-xboard-client/internal/config"
	"github.com/your-username/singbox-xboard-client/pkg/xboard"
)

// Manager 订阅管理器
type Manager struct {
	config      *config.Config
	client      *xboard.Client
	cron        *cron.Cron
	logger      *logrus.Logger
	mu          sync.RWMutex
	lastUpdate  time.Time
	lastConfig  map[string]interface{}
	updateHooks []func(map[string]interface{})
}

// NewManager 创建订阅管理器
func NewManager() *Manager {
	return &Manager{
		logger:      logrus.New(),
		updateHooks: make([]func(map[string]interface{}), 0),
	}
}

// Initialize 初始化订阅管理器
func (m *Manager) Initialize(cfg *config.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = cfg

	// 如果有订阅 URL，创建客户端
	if cfg.Subscription.URL != "" {
		baseURL, token, err := xboard.ParseSubscriptionURL(cfg.Subscription.URL)
		if err != nil {
			// 尝试直接使用 URL 和 token
			if cfg.Subscription.Token != "" {
				m.client = xboard.NewClient(cfg.Subscription.URL, cfg.Subscription.Token)
			} else {
				return fmt.Errorf("解析订阅 URL 失败: %w", err)
			}
		} else {
			m.client = xboard.NewClient(baseURL, token)
		}
		m.client.SetLogger(m.logger)
	}

	// 设置自动更新
	if cfg.Subscription.AutoUpdate && cfg.Subscription.UpdateInterval > 0 {
		m.setupAutoUpdate()
	}

	return nil
}

// UpdateSubscription 更新订阅
func (m *Manager) UpdateSubscription(url string) error {
	m.logger.Info("开始更新订阅")

	// 解析订阅 URL
	baseURL, token, err := xboard.ParseSubscriptionURL(url)
	if err != nil {
		return fmt.Errorf("解析订阅 URL 失败: %w", err)
	}

	// 创建临时客户端
	client := xboard.NewClient(baseURL, token)
	client.SetLogger(m.logger)

	// 获取配置
	singboxConfig, err := client.GetSingboxConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	// 保存配置
	if err := m.saveConfig(singboxConfig); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	// 更新状态
	m.mu.Lock()
	m.lastUpdate = time.Now()
	m.lastConfig = singboxConfig
	m.client = client
	
	// 更新应用配置
	if m.config != nil {
		m.config.Subscription.URL = url
		m.config.Subscription.Token = token
	}
	m.mu.Unlock()

	// 触发更新钩子
	m.notifyUpdate(singboxConfig)

	m.logger.Info("订阅更新成功")
	return nil
}

// GetLastConfig 获取最后的配置
func (m *Manager) GetLastConfig() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastConfig
}

// GetLastUpdate 获取最后更新时间
func (m *Manager) GetLastUpdate() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastUpdate
}

// RefreshSubscription 刷新当前订阅
func (m *Manager) RefreshSubscription() error {
	m.mu.RLock()
	client := m.client
	m.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("未配置订阅")
	}

	m.logger.Info("刷新订阅")

	// 获取配置
	singboxConfig, err := client.GetSingboxConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	// 保存配置
	if err := m.saveConfig(singboxConfig); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	// 更新状态
	m.mu.Lock()
	m.lastUpdate = time.Now()
	m.lastConfig = singboxConfig
	m.mu.Unlock()

	// 触发更新钩子
	m.notifyUpdate(singboxConfig)

	m.logger.Info("订阅刷新成功")
	return nil
}

// GetUserInfo 获取用户信息
func (m *Manager) GetUserInfo() (*xboard.UserInfo, error) {
	m.mu.RLock()
	client := m.client
	m.mu.RUnlock()

	if client == nil {
		return nil, fmt.Errorf("未配置订阅")
	}

	return client.GetUserInfo()
}

// GetNodeList 获取节点列表
func (m *Manager) GetNodeList() ([]xboard.NodeInfo, error) {
	m.mu.RLock()
	client := m.client
	m.mu.RUnlock()

	if client == nil {
		return nil, fmt.Errorf("未配置订阅")
	}

	return client.GetNodeList()
}

// ReportTraffic 上报流量
func (m *Manager) ReportTraffic(upload, download int64, nodeID int) error {
	m.mu.RLock()
	client := m.client
	m.mu.RUnlock()

	if client == nil {
		return fmt.Errorf("未配置订阅")
	}

	return client.ReportTraffic(upload, download, nodeID)
}

// OnUpdate 注册更新钩子
func (m *Manager) OnUpdate(hook func(map[string]interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateHooks = append(m.updateHooks, hook)
}

// setupAutoUpdate 设置自动更新
func (m *Manager) setupAutoUpdate() {
	if m.cron != nil {
		m.cron.Stop()
	}

	m.cron = cron.New()
	
	// 添加定时任务
	spec := fmt.Sprintf("@every %dm", m.config.Subscription.UpdateInterval)
	m.cron.AddFunc(spec, func() {
		m.logger.Info("执行自动更新")
		if err := m.RefreshSubscription(); err != nil {
			m.logger.Errorf("自动更新失败: %v", err)
		}
	})

	m.cron.Start()
	m.logger.Infof("已启用自动更新，间隔: %d 分钟", m.config.Subscription.UpdateInterval)
}

// saveConfig 保存配置到文件
func (m *Manager) saveConfig(singboxConfig map[string]interface{}) error {
	// 获取配置目录
	configDir := config.GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 保存为 JSON 文件
	configPath := filepath.Join(configDir, "singbox.json")
	data, err := json.MarshalIndent(singboxConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	m.logger.Debugf("配置已保存到: %s", configPath)
	return nil
}

// notifyUpdate 通知更新
func (m *Manager) notifyUpdate(singboxConfig map[string]interface{}) {
	m.mu.RLock()
	hooks := make([]func(map[string]interface{}), len(m.updateHooks))
	copy(hooks, m.updateHooks)
	m.mu.RUnlock()

	for _, hook := range hooks {
		go func(h func(map[string]interface{})) {
			defer func() {
				if r := recover(); r != nil {
					m.logger.Errorf("更新钩子执行失败: %v", r)
				}
			}()
			h(singboxConfig)
		}(hook)
	}
}

// LoadCachedConfig 加载缓存的配置
func (m *Manager) LoadCachedConfig() (map[string]interface{}, error) {
	configPath := filepath.Join(config.GetConfigDir(), "singbox.json")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取缓存配置失败: %w", err)
	}

	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析缓存配置失败: %w", err)
	}

	m.mu.Lock()
	m.lastConfig = cfg
	m.mu.Unlock()

	return cfg, nil
}

// Stop 停止订阅管理器
func (m *Manager) Stop() {
	if m.cron != nil {
		m.cron.Stop()
		m.logger.Info("已停止自动更新")
	}
}

// SetLogger 设置日志记录器
func (m *Manager) SetLogger(logger *logrus.Logger) {
	m.logger = logger
}