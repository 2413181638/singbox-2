package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"singbox-app/internal/config"
	"singbox-app/internal/proxy"
)

func SetupRoutes(router *gin.Engine, proxyService *proxy.Service, cfg *config.Config) {
	// 静态文件
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "./web/static")
	
	// 主页
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "SingBox 管理面板",
		})
	})
	
	// API路由组
	api := router.Group("/api")
	{
		// 服务状态
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"running": proxyService.IsRunning(),
				"stats":   proxyService.GetStats(),
			})
		})
		
		// 获取配置
		api.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, cfg)
		})
		
		// 更新配置
		api.POST("/config", func(c *gin.Context) {
			var newConfig config.Config
			if err := c.ShouldBindJSON(&newConfig); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			
			// 这里应该保存配置并重启服务
			c.JSON(http.StatusOK, gin.H{"message": "配置已更新"})
		})
		
		// 启动服务
		api.POST("/start", func(c *gin.Context) {
			if proxyService.IsRunning() {
				c.JSON(http.StatusBadRequest, gin.H{"error": "服务已在运行"})
				return
			}
			
			if err := proxyService.Start(c.Request.Context()); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			c.JSON(http.StatusOK, gin.H{"message": "服务已启动"})
		})
		
		// 停止服务
		api.POST("/stop", func(c *gin.Context) {
			if !proxyService.IsRunning() {
				c.JSON(http.StatusBadRequest, gin.H{"error": "服务未运行"})
				return
			}
			
			if err := proxyService.Stop(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			c.JSON(http.StatusOK, gin.H{"message": "服务已停止"})
		})
	}
}