package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/your-username/singbox-xboard-client/internal/config"
	"github.com/your-username/singbox-xboard-client/internal/singbox"
	"github.com/your-username/singbox-xboard-client/internal/subscription"
	"github.com/your-username/singbox-xboard-client/internal/ui"
)

var (
	// 版本信息，编译时注入
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "singbox-xboard",
	Short: "基于 sing-box 的 xboard 客户端",
	Long:  `一个基于 sing-box 内核的跨平台客户端，支持 xboard 面板订阅`,
	Run: func(cmd *cobra.Command, args []string) {
		// 默认启动 GUI 模式
		runGUI()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Singbox Xboard Client\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行客户端",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		mode, _ := cmd.Flags().GetString("mode")

		// 加载配置
		cfg, err := config.Load(configFile)
		if err != nil {
			logrus.Fatalf("加载配置失败: %v", err)
		}

		switch mode {
		case "cli":
			runCLI(cfg)
		case "gui":
			runGUI()
		default:
			runGUI()
		}
	},
}

var subscribeCmd = &cobra.Command{
	Use:   "subscribe [url]",
	Short: "更新订阅",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		
		// 创建订阅管理器
		subMgr := subscription.NewManager()
		
		// 更新订阅
		if err := subMgr.UpdateSubscription(url); err != nil {
			logrus.Fatalf("更新订阅失败: %v", err)
		}
		
		fmt.Println("订阅更新成功")
	},
}

func init() {
	// 设置日志
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 添加子命令
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(subscribeCmd)

	// 设置运行命令的标志
	runCmd.Flags().StringP("config", "c", "", "配置文件路径")
	runCmd.Flags().StringP("mode", "m", "gui", "运行模式 (gui/cli)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runGUI 启动图形界面
func runGUI() {
	logrus.Info("启动 GUI 模式")
	
	// 创建 UI 服务器
	server := ui.NewServer()
	
	// 启动 UI
	if err := server.Start(); err != nil {
		logrus.Fatalf("启动 UI 失败: %v", err)
	}
}

// runCLI 启动命令行模式
func runCLI(cfg *config.Config) {
	logrus.Info("启动 CLI 模式")
	
	// 创建 sing-box 管理器
	manager := singbox.NewManager(cfg)
	
	// 启动 sing-box
	if err := manager.Start(); err != nil {
		logrus.Fatalf("启动 sing-box 失败: %v", err)
	}
	
	// 等待退出信号
	select {}
}