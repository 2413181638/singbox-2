package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"singbox-app/cmd"
	"singbox-app/internal/logger"
)

func main() {
	// 初始化日志
	logger.Init()

	// 创建上下文，用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n正在关闭应用...")
		cancel()
	}()

	// 执行命令
	if err := cmd.Execute(ctx); err != nil {
		logger.Error("应用启动失败: %v", err)
		os.Exit(1)
	}
}