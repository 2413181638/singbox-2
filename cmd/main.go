package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"singbox-xboard-client/internal/app"
	"singbox-xboard-client/internal/config"
	"singbox-xboard-client/internal/core"
	"singbox-xboard-client/internal/database"
	"singbox-xboard-client/internal/xboard"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	db, err := database.Init(cfg.DatabasePath)
	if err != nil {
		fmt.Printf("Failed to init database: %v\n", err)
		os.Exit(1)
	}

	// 初始化核心
	core := core.New(cfg)

	// 初始化 XBoard 客户端
	xboardClient := xboard.New(cfg.XBoard)

	// 创建应用实例
	app := app.New(cfg, db, core, xboardClient)

	// 创建 Wails 应用
	err = wails.Run(&options.App{
		Title:  "SingBox XBoard Client",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnDomReady:       app.DomReady,
		OnShutdown:       app.Shutdown,
		Frameless:        false,
		MinWidth:         800,
		MinHeight:        600,
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}