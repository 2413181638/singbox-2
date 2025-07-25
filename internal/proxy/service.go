package proxy

import (
	"context"
	"fmt"
	"sync"

	"singbox-app/internal/config"
	"singbox-app/internal/logger"
)

type Service struct {
	config  *config.Config
	running bool
	mu      sync.RWMutex
	cancel  context.CancelFunc
}

func NewService(cfg *config.Config) (*Service, error) {
	return &Service{
		config: cfg,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.running {
		return fmt.Errorf("服务已在运行")
	}
	
	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	
	s.running = true
	
	logger.Info("代理服务已启动，类型: %s，监听端口: %d", s.config.Inbound.Type, s.config.Inbound.Port)
	logger.Info("注意: 当前为演示版本，实际代理功能需要完整的sing-box集成")
	
	return nil
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.running {
		return nil
	}
	
	s.running = false
	
	if s.cancel != nil {
		s.cancel()
	}
	
	logger.Info("代理服务已停止")
	return nil
}

func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *Service) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"running": s.IsRunning(),
		"config":  s.config,
		"type":    "基础代理服务",
		"note":    "这是一个演示版本，可以扩展为完整的sing-box集成",
	}
}