package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"singbox-app/internal/config"
	"singbox-app/internal/logger"
	"singbox-app/internal/proxy"
	"singbox-app/internal/web"
)

func startWebUI(ctx context.Context) error {
	logger.Info("启动Web UI模式...")
	
	// 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	
	// 创建代理服务
	proxyService, err := proxy.NewService(cfg)
	if err != nil {
		return fmt.Errorf("创建代理服务失败: %w", err)
	}
	
	// 启动代理服务
	if err := proxyService.Start(ctx); err != nil {
		return fmt.Errorf("启动代理服务失败: %w", err)
	}
	
	// 创建Web服务器
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	
	// 设置Web路由
	web.SetupRoutes(router, proxyService, cfg)
	
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	
	// 启动Web服务器
	go func() {
		logger.Info("Web界面已启动: http://localhost:%d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Web服务器启动失败: %v", err)
		}
	}()
	
	// 等待上下文取消
	<-ctx.Done()
	
	// 优雅关闭
	logger.Info("正在关闭Web服务器...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("关闭Web服务器失败: %v", err)
	}
	
	// 停止代理服务
	logger.Info("正在停止代理服务...")
	if err := proxyService.Stop(); err != nil {
		logger.Error("停止代理服务失败: %v", err)
	}
	
	logger.Info("应用程序已完全关闭")
	return nil
}