package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"singbox-xboard-client/internal/config"
	"singbox-xboard-client/internal/core"
	"singbox-xboard-client/internal/database"
	"singbox-xboard-client/internal/xboard"
)

type App struct {
	config       *config.Config
	db           *database.Database
	core         *core.Core
	xboardClient *xboard.Client
	ctx          context.Context
	cancel       context.CancelFunc
}

type ServerInfo struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	ServerType string `json:"server_type"`
	IsActive   bool   `json:"is_active"`
	Ping       int    `json:"ping"`
}

type UserStatus struct {
	Email     string `json:"email"`
	Upload    int64  `json:"upload"`
	Download  int64  `json:"download"`
	Total     int64  `json:"total"`
	Remaining int64  `json:"remaining"`
	ExpiredAt int64  `json:"expired_at"`
	IsActive  bool   `json:"is_active"`
}

type ConnectionStatus struct {
	IsRunning bool                   `json:"is_running"`
	Stats     map[string]interface{} `json:"stats"`
}

func New(cfg *config.Config, db *database.Database, core *core.Core, xboardClient *xboard.Client) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		config:       cfg,
		db:           db,
		core:         core,
		xboardClient: xboardClient,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (a *App) Startup(ctx context.Context) {
	log.Println("应用启动中...")
	
	// 启动定时任务
	go a.startScheduledTasks()
}

func (a *App) DomReady(ctx context.Context) {
	log.Println("前端界面已就绪")
}

func (a *App) Shutdown(ctx context.Context) {
	log.Println("应用关闭中...")
	a.cancel()
	if err := a.core.Stop(); err != nil {
		log.Printf("停止核心失败: %v", err)
	}
}

// 登录到 XBoard
func (a *App) Login(email, password string) error {
	token, err := a.xboardClient.Login(email, password)
	if err != nil {
		return err
	}

	// 保存 token 到配置
	a.config.XBoard.Token = token
	if err := a.config.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	// 同步订阅
	return a.SyncSubscription()
}

// 同步订阅
func (a *App) SyncSubscription() error {
	if a.config.XBoard.Token == "" {
		return fmt.Errorf("未登录")
	}

	nodes, err := a.xboardClient.GetSubscription(a.config.XBoard.Token)
	if err != nil {
		return err
	}

	// 清空现有服务器
	if err := a.db.Exec("DELETE FROM servers").Error; err != nil {
		return err
	}

	// 添加新服务器
	for _, node := range nodes {
		server := &database.Server{
			Name:       node.Name,
			Host:       node.Host,
			Port:       node.Port,
			ServerType: node.Type,
			Settings:   node.Settings,
			Tags:       node.Tags,
			IsActive:   true,
		}

		if err := a.db.CreateServer(server); err != nil {
			log.Printf("添加服务器失败: %v", err)
			continue
		}
	}

	// 更新核心配置
	servers, err := a.db.GetActiveServers()
	if err != nil {
		return err
	}

	return a.core.UpdateConfig(servers)
}

// 获取服务器列表
func (a *App) GetServers() ([]ServerInfo, error) {
	servers, err := a.db.GetServers()
	if err != nil {
		return nil, err
	}

	var serverInfos []ServerInfo
	for _, server := range servers {
		serverInfos = append(serverInfos, ServerInfo{
			ID:         server.ID,
			Name:       server.Name,
			Host:       server.Host,
			Port:       server.Port,
			ServerType: server.ServerType,
			IsActive:   server.IsActive,
			Ping:       server.Ping,
		})
	}

	return serverInfos, nil
}

// 切换服务器状态
func (a *App) ToggleServer(id uint) error {
	var server database.Server
	if err := a.db.First(&server, id).Error; err != nil {
		return err
	}

	server.IsActive = !server.IsActive
	if err := a.db.UpdateServer(&server); err != nil {
		return err
	}

	// 更新核心配置
	servers, err := a.db.GetActiveServers()
	if err != nil {
		return err
	}

	return a.core.UpdateConfig(servers)
}

// 启动连接
func (a *App) StartConnection() error {
	return a.core.Start()
}

// 停止连接
func (a *App) StopConnection() error {
	return a.core.Stop()
}

// 获取连接状态
func (a *App) GetConnectionStatus() (*ConnectionStatus, error) {
	status := &ConnectionStatus{
		IsRunning: a.core.IsRunning(),
	}

	if status.IsRunning {
		stats, err := a.core.GetStats()
		if err != nil {
			log.Printf("获取统计信息失败: %v", err)
		} else {
			status.Stats = stats
		}
	}

	return status, nil
}

// 获取用户状态
func (a *App) GetUserStatus() (*UserStatus, error) {
	if a.config.XBoard.Token == "" {
		return nil, fmt.Errorf("未登录")
	}

	userInfo, err := a.xboardClient.GetUserInfo(a.config.XBoard.Token)
	if err != nil {
		return nil, err
	}

	return &UserStatus{
		Email:     userInfo.Email,
		Upload:    userInfo.Transfer.Up,
		Download:  userInfo.Transfer.Down,
		Total:     userInfo.Transfer.Total,
		Remaining: userInfo.Transfer.Remaining,
		ExpiredAt: userInfo.ExpiredAt,
		IsActive:  userInfo.Status == 1,
	}, nil
}

// 测试服务器延迟
func (a *App) PingServer(id uint) (int, error) {
	var server database.Server
	if err := a.db.First(&server, id).Error; err != nil {
		return 0, err
	}

	// 简单的 TCP 连接测试
	start := time.Now()
	// TODO: 实现真正的 ping 测试
	ping := int(time.Since(start).Milliseconds())
	
	server.Ping = ping
	if err := a.db.UpdateServer(&server); err != nil {
		return 0, err
	}

	return ping, nil
}

// 获取连接日志
func (a *App) GetConnectionLogs() ([]database.ConnectionLog, error) {
	return a.db.GetConnectionLogs(100)
}

// 定时任务
func (a *App) startScheduledTasks() {
	// 定时同步订阅
	syncTicker := time.NewTicker(time.Duration(a.config.XBoard.Interval) * time.Second)
	defer syncTicker.Stop()

	// 定时上报流量
	trafficTicker := time.NewTicker(60 * time.Second) // 每分钟上报一次
	defer trafficTicker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-syncTicker.C:
			if err := a.SyncSubscription(); err != nil {
				log.Printf("定时同步订阅失败: %v", err)
			}
		case <-trafficTicker.C:
			if err := a.reportTraffic(); err != nil {
				log.Printf("上报流量失败: %v", err)
			}
		}
	}
}

// 上报流量
func (a *App) reportTraffic() error {
	if a.config.XBoard.Token == "" || !a.core.IsRunning() {
		return nil
	}

	stats, err := a.core.GetStats()
	if err != nil {
		return err
	}

	// 提取上传和下载流量
	var upload, download int64
	if upValue, ok := stats["up"]; ok {
		if up, ok := upValue.(float64); ok {
			upload = int64(up)
		}
	}
	if downValue, ok := stats["down"]; ok {
		if down, ok := downValue.(float64); ok {
			download = int64(down)
		}
	}

	return a.xboardClient.ReportTraffic(a.config.XBoard.Token, upload, download)
}

// 获取配置
func (a *App) GetConfig() *config.Config {
	return a.config
}

// 更新配置
func (a *App) UpdateConfig(newConfig *config.Config) error {
	a.config = newConfig
	return a.config.Save()
}