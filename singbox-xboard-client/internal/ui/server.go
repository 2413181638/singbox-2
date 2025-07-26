package ui

import (
	"embed"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/your-username/singbox-xboard-client/internal/config"
	"github.com/your-username/singbox-xboard-client/internal/singbox"
	"github.com/your-username/singbox-xboard-client/internal/subscription"
)

//go:embed static/*
var staticFiles embed.FS

// Server Web UI 服务器
type Server struct {
	config      *config.Config
	engine      *gin.Engine
	subManager  *subscription.Manager
	sbManager   *singbox.Manager
	logger      *logrus.Logger
	upgrader    websocket.Upgrader
}

// NewServer 创建 UI 服务器
func NewServer() *Server {
	// 设置 Gin 为发布模式
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		engine: gin.New(),
		logger: logrus.New(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源，生产环境应该限制
			},
		},
	}

	// 初始化路由
	s.setupRoutes()

	return s
}

// Start 启动服务器
func (s *Server) Start() error {
	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	s.config = cfg

	// 创建管理器
	s.subManager = subscription.NewManager()
	s.subManager.SetLogger(s.logger)
	
	s.sbManager = singbox.NewManager(cfg)
	s.sbManager.SetLogger(s.logger)

	// 初始化订阅管理器
	if err := s.subManager.Initialize(cfg); err != nil {
		s.logger.Warnf("初始化订阅管理器失败: %v", err)
	}

	// 注册订阅更新钩子
	s.subManager.OnUpdate(func(config map[string]interface{}) {
		s.logger.Info("收到订阅更新，重新加载配置")
		if err := s.sbManager.UpdateConfig(config); err != nil {
			s.logger.Errorf("更新 sing-box 配置失败: %v", err)
		}
	})

	// 尝试加载缓存的配置
	if cachedConfig, err := s.subManager.LoadCachedConfig(); err == nil {
		s.logger.Info("加载缓存配置")
		s.sbManager.UpdateConfig(cachedConfig)
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.UI.Listen, cfg.UI.Port)
	s.logger.Infof("启动 Web UI 服务器: http://%s", addr)
	
	// 在浏览器中打开
	go func() {
		time.Sleep(time.Second)
		openBrowser(fmt.Sprintf("http://localhost:%d", cfg.UI.Port))
	}()

	return s.engine.Run(addr)
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 静态文件
	s.engine.StaticFS("/static", http.FS(staticFiles))
	
	// API 路由
	api := s.engine.Group("/api")
	{
		// 系统信息
		api.GET("/status", s.handleGetStatus)
		api.GET("/config", s.handleGetConfig)
		api.POST("/config", s.handleUpdateConfig)
		
		// 订阅管理
		api.GET("/subscription", s.handleGetSubscription)
		api.POST("/subscription", s.handleUpdateSubscription)
		api.POST("/subscription/refresh", s.handleRefreshSubscription)
		
		// 节点管理
		api.GET("/nodes", s.handleGetNodes)
		api.POST("/node/select", s.handleSelectNode)
		
		// Sing-box 控制
		api.POST("/singbox/start", s.handleStartSingbox)
		api.POST("/singbox/stop", s.handleStopSingbox)
		api.POST("/singbox/restart", s.handleRestartSingbox)
		
		// WebSocket
		api.GET("/ws", s.handleWebSocket)
	}
	
	// 主页
	s.engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})
}

// handleGetStatus 获取状态
func (s *Server) handleGetStatus(c *gin.Context) {
	upload, download, uptime := s.sbManager.GetStats()
	
	status := gin.H{
		"running": s.sbManager.IsRunning(),
		"uptime":  uptime.Seconds(),
		"stats": gin.H{
			"upload":   upload,
			"download": download,
		},
	}
	
	// 获取用户信息
	if userInfo, err := s.subManager.GetUserInfo(); err == nil {
		status["user"] = userInfo
	}
	
	// 获取最后更新时间
	if lastUpdate := s.subManager.GetLastUpdate(); !lastUpdate.IsZero() {
		status["lastUpdate"] = lastUpdate
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// handleGetConfig 获取配置
func (s *Server) handleGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    s.config,
	})
}

// handleUpdateConfig 更新配置
func (s *Server) handleUpdateConfig(c *gin.Context) {
	var newConfig config.Config
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	// 保存配置
	if err := config.Save(&newConfig, config.GetDefaultConfigPath()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	s.config = &newConfig
	
	// 重新初始化订阅管理器
	if err := s.subManager.Initialize(&newConfig); err != nil {
		s.logger.Warnf("重新初始化订阅管理器失败: %v", err)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// handleGetSubscription 获取订阅信息
func (s *Server) handleGetSubscription(c *gin.Context) {
	data := gin.H{
		"url":        s.config.Subscription.URL,
		"lastUpdate": s.subManager.GetLastUpdate(),
	}
	
	// 获取节点列表
	if nodes, err := s.subManager.GetNodeList(); err == nil {
		data["nodes"] = nodes
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// handleUpdateSubscription 更新订阅
func (s *Server) handleUpdateSubscription(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	// 更新订阅
	if err := s.subManager.UpdateSubscription(req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	// 保存配置
	s.config.Subscription.URL = req.URL
	config.Save(s.config, config.GetDefaultConfigPath())
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// handleRefreshSubscription 刷新订阅
func (s *Server) handleRefreshSubscription(c *gin.Context) {
	if err := s.subManager.RefreshSubscription(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// handleGetNodes 获取节点列表
func (s *Server) handleGetNodes(c *gin.Context) {
	nodes, err := s.subManager.GetNodeList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nodes,
	})
}

// handleSelectNode 选择节点
func (s *Server) handleSelectNode(c *gin.Context) {
	var req struct {
		NodeID int `json:"node_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	// TODO: 实现节点选择逻辑
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// handleStartSingbox 启动 sing-box
func (s *Server) handleStartSingbox(c *gin.Context) {
	if err := s.sbManager.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// handleStopSingbox 停止 sing-box
func (s *Server) handleStopSingbox(c *gin.Context) {
	if err := s.sbManager.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// handleRestartSingbox 重启 sing-box
func (s *Server) handleRestartSingbox(c *gin.Context) {
	if err := s.sbManager.Restart(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// handleWebSocket 处理 WebSocket 连接
func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.logger.Errorf("WebSocket 升级失败: %v", err)
		return
	}
	defer conn.Close()
	
	// 定期发送状态更新
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			upload, download, uptime := s.sbManager.GetStats()
			status := map[string]interface{}{
				"type":    "status",
				"running": s.sbManager.IsRunning(),
				"uptime":  uptime.Seconds(),
				"stats": map[string]interface{}{
					"upload":   upload,
					"download": download,
				},
			}
			
			if err := conn.WriteJSON(status); err != nil {
				s.logger.Debugf("WebSocket 写入失败: %v", err)
				return
			}
		}
	}
}

// openBrowser 在浏览器中打开 URL
func openBrowser(url string) {
	var err error
	
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	
	if err != nil {
		logrus.Debugf("无法打开浏览器: %v", err)
	}
}