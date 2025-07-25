package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	configFile string
	webUI      bool
	port       int
)

var rootCmd = &cobra.Command{
	Use:   "singbox-app",
	Short: "基于sing-box的代理应用程序",
	Long:  `一个功能完整的基于sing-box的代理应用程序，支持多种协议和配置方式。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApp(cmd.Context())
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&webUI, "web", "w", false, "启用Web管理界面")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Web界面端口")
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func runApp(ctx context.Context) error {
	fmt.Println("启动 SingBox 应用程序...")
	
	if webUI {
		return startWebUI(ctx)
	}
	
	return startCLI(ctx)
}