package cmd

import (
	"context"
	"fmt"

	"singbox-app/internal/config"
	"singbox-app/internal/logger"
	"singbox-app/internal/proxy"
)

func startCLI(ctx context.Context) error {
	logger.Info("启动CLI模式...")
	
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
	
	// 启动服务
	if err := proxyService.Start(ctx); err != nil {
		return fmt.Errorf("启动代理服务失败: %w", err)
	}
	
	logger.Info("代理服务已启动，监听端口: %d", cfg.Inbound.Port)
	logger.Info("按 Ctrl+C 停止服务")
	
	// 等待上下文取消
	<-ctx.Done()
	
	// 停止服务
	logger.Info("正在停止代理服务...")
	if err := proxyService.Stop(); err != nil {
		logger.Error("停止代理服务失败: %v", err)
	}
	
	logger.Info("代理服务已停止")
	return nil
}