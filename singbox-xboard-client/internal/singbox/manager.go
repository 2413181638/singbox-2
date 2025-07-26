package singbox

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/your-username/singbox-xboard-client/internal/config"
)

// Manager sing-box 管理器
type Manager struct {
	config       *config.Config
	process      *exec.Cmd
	ctx          context.Context
	cancel       context.CancelFunc
	logger       *logrus.Logger
	mu           sync.Mutex
	isRunning    bool
	configPath   string
	statsTracker *StatsTracker
}

// StatsTracker 流量统计跟踪器
type StatsTracker struct {
	mu         sync.RWMutex
	upload     int64
	download   int64
	startTime  time.Time
	lastUpdate time.Time
}

// NewManager 创建 sing-box 管理器
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:       cfg,
		logger:       logrus.New(),
		statsTracker: &StatsTracker{},
	}
}

// Start 启动 sing-box
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("sing-box 已在运行")
	}

	// 创建上下文
	m.ctx, m.cancel = context.WithCancel(context.Background())

	// 准备配置文件
	if err := m.prepareConfig(); err != nil {
		return fmt.Errorf("准备配置失败: %w", err)
	}

	// 查找 sing-box 可执行文件
	singboxPath, err := m.findSingboxBinary()
	if err != nil {
		return fmt.Errorf("找不到 sing-box: %w", err)
	}

	// 创建命令
	m.process = exec.CommandContext(m.ctx, singboxPath, "run", "-c", m.configPath)
	
	// 设置输出
	m.process.Stdout = m.logger.Writer()
	m.process.Stderr = m.logger.Writer()

	// 启动进程
	if err := m.process.Start(); err != nil {
		return fmt.Errorf("启动 sing-box 失败: %w", err)
	}

	m.isRunning = true
	m.statsTracker.startTime = time.Now()
	m.statsTracker.lastUpdate = time.Now()

	// 监控进程
	go m.monitorProcess()

	m.logger.Info("sing-box 已启动")
	return nil
}

// Stop 停止 sing-box
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("sing-box 未运行")
	}

	// 取消上下文
	if m.cancel != nil {
		m.cancel()
	}

	// 等待进程退出
	if m.process != nil {
		// 给进程一些时间优雅退出
		done := make(chan error, 1)
		go func() {
			done <- m.process.Wait()
		}()

		select {
		case <-done:
			// 进程已退出
		case <-time.After(5 * time.Second):
			// 强制终止
			if err := m.process.Process.Kill(); err != nil {
				m.logger.Warnf("强制终止失败: %v", err)
			}
		}
	}

	m.isRunning = false
	m.logger.Info("sing-box 已停止")
	return nil
}

// Restart 重启 sing-box
func (m *Manager) Restart() error {
	if m.isRunning {
		if err := m.Stop(); err != nil {
			return fmt.Errorf("停止失败: %w", err)
		}
		// 等待一下
		time.Sleep(time.Second)
	}
	return m.Start()
}

// IsRunning 检查是否正在运行
func (m *Manager) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isRunning
}

// UpdateConfig 更新配置
func (m *Manager) UpdateConfig(singboxConfig map[string]interface{}) error {
	// 保存新配置
	configPath := filepath.Join(config.GetConfigDir(), "singbox.json")
	data, err := json.MarshalIndent(singboxConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	// 如果正在运行，重启以应用新配置
	if m.IsRunning() {
		m.logger.Info("重启 sing-box 以应用新配置")
		return m.Restart()
	}

	return nil
}

// GetStats 获取统计信息
func (m *Manager) GetStats() (upload, download int64, uptime time.Duration) {
	m.statsTracker.mu.RLock()
	defer m.statsTracker.mu.RUnlock()

	upload = m.statsTracker.upload
	download = m.statsTracker.download
	if m.isRunning {
		uptime = time.Since(m.statsTracker.startTime)
	}
	return
}

// UpdateStats 更新统计信息
func (m *Manager) UpdateStats(upload, download int64) {
	m.statsTracker.mu.Lock()
	defer m.statsTracker.mu.Unlock()

	m.statsTracker.upload += upload
	m.statsTracker.download += download
	m.statsTracker.lastUpdate = time.Now()
}

// prepareConfig 准备配置文件
func (m *Manager) prepareConfig() error {
	// 获取配置目录
	configDir := config.GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 检查是否有缓存的 sing-box 配置
	m.configPath = filepath.Join(configDir, "singbox.json")
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// 创建默认配置
		defaultConfig := m.createDefaultConfig()
		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化默认配置失败: %w", err)
		}

		if err := os.WriteFile(m.configPath, data, 0644); err != nil {
			return fmt.Errorf("写入默认配置失败: %w", err)
		}
	}

	return nil
}

// findSingboxBinary 查找 sing-box 可执行文件
func (m *Manager) findSingboxBinary() (string, error) {
	// 优先查找系统路径
	if path, err := exec.LookPath("sing-box"); err == nil {
		return path, nil
	}

	// 查找应用目录
	exePath, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exePath)
		// 同目录
		singboxPath := filepath.Join(dir, "sing-box")
		if _, err := os.Stat(singboxPath); err == nil {
			return singboxPath, nil
		}
		// bin 子目录
		singboxPath = filepath.Join(dir, "bin", "sing-box")
		if _, err := os.Stat(singboxPath); err == nil {
			return singboxPath, nil
		}
	}

	// 查找配置目录
	configDir := config.GetConfigDir()
	singboxPath := filepath.Join(configDir, "sing-box")
	if _, err := os.Stat(singboxPath); err == nil {
		return singboxPath, nil
	}

	return "", fmt.Errorf("找不到 sing-box 可执行文件")
}

// monitorProcess 监控进程
func (m *Manager) monitorProcess() {
	err := m.process.Wait()
	
	m.mu.Lock()
	m.isRunning = false
	m.mu.Unlock()

	if err != nil {
		m.logger.Errorf("sing-box 进程退出: %v", err)
	} else {
		m.logger.Info("sing-box 进程正常退出")
	}
}

// createDefaultConfig 创建默认配置
func (m *Manager) createDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"log": map[string]interface{}{
			"level":     "info",
			"timestamp": true,
		},
		"dns": map[string]interface{}{
			"servers": []map[string]interface{}{
				{
					"tag":     "remote",
					"address": "https://1.1.1.1/dns-query",
					"detour":  "direct",
				},
				{
					"tag":     "local",
					"address": "https://223.5.5.5/dns-query",
					"detour":  "direct",
				},
			},
			"rules": []map[string]interface{}{
				{
					"geosite": []string{"cn"},
					"server":  "local",
				},
			},
			"final": "remote",
		},
		"inbounds": []map[string]interface{}{
			{
				"type":        "mixed",
				"tag":         "mixed-in",
				"listen":      "127.0.0.1",
				"listen_port": 7890,
				"sniff":       true,
			},
		},
		"outbounds": []map[string]interface{}{
			{
				"type": "direct",
				"tag":  "direct",
			},
			{
				"type": "block",
				"tag":  "block",
			},
			{
				"type": "dns",
				"tag":  "dns-out",
			},
		},
		"route": map[string]interface{}{
			"rules": []map[string]interface{}{
				{
					"protocol": "dns",
					"outbound": "dns-out",
				},
			},
			"final": "direct",
		},
	}
}

// SetLogger 设置日志记录器
func (m *Manager) SetLogger(logger *logrus.Logger) {
	m.logger = logger
}