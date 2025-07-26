package xboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Client xboard API 客户端
type Client struct {
	baseURL    string
	token      string
	httpClient *resty.Client
	logger     *logrus.Logger
}

// NewClient 创建新的 xboard 客户端
func NewClient(baseURL, token string) *Client {
	// 处理 baseURL，确保格式正确
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}

	client := &Client{
		baseURL: baseURL,
		token:   token,
		logger:  logrus.New(),
	}

	// 创建 HTTP 客户端
	client.httpClient = resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetHeader("User-Agent", "SingboxXboardClient/1.0").
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			// 添加认证信息
			if client.token != "" {
				req.SetHeader("Authorization", "Bearer "+client.token)
			}
			return nil
		})

	return client
}

// GetSubscription 获取订阅信息
func (c *Client) GetSubscription() (*SubscriptionResponse, error) {
	c.logger.Debug("获取订阅信息")

	resp, err := c.httpClient.R().
		SetResult(&SubscriptionResponse{}).
		Get(c.baseURL + "/api/v1/client/subscribe")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("服务器返回错误: %s", resp.Status())
	}

	result := resp.Result().(*SubscriptionResponse)
	return result, nil
}

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo() (*UserInfo, error) {
	c.logger.Debug("获取用户信息")

	resp, err := c.httpClient.R().
		SetResult(&UserInfo{}).
		Get(c.baseURL + "/api/v1/user/info")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("服务器返回错误: %s", resp.Status())
	}

	result := resp.Result().(*UserInfo)
	return result, nil
}

// GetNodeList 获取节点列表
func (c *Client) GetNodeList() ([]NodeInfo, error) {
	c.logger.Debug("获取节点列表")

	var nodes []NodeInfo
	resp, err := c.httpClient.R().
		SetResult(&nodes).
		Get(c.baseURL + "/api/v1/user/node")

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("服务器返回错误: %s", resp.Status())
	}

	return nodes, nil
}

// ReportTraffic 上报流量使用情况
func (c *Client) ReportTraffic(upload, download int64, nodeID int) error {
	c.logger.Debugf("上报流量: 上传=%d, 下载=%d, 节点=%d", upload, download, nodeID)

	data := map[string]interface{}{
		"upload":   upload,
		"download": download,
		"node_id":  nodeID,
	}

	resp, err := c.httpClient.R().
		SetBody(data).
		Post(c.baseURL + "/api/v1/user/traffic")

	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("服务器返回错误: %s", resp.Status())
	}

	return nil
}

// GetSingboxConfig 获取 sing-box 格式的配置
func (c *Client) GetSingboxConfig() (map[string]interface{}, error) {
	// 获取订阅信息
	sub, err := c.GetSubscription()
	if err != nil {
		return nil, fmt.Errorf("获取订阅失败: %w", err)
	}

	// 构建 sing-box 配置
	config := map[string]interface{}{
		"log": map[string]interface{}{
			"level":     "info",
			"timestamp": true,
		},
		"dns": map[string]interface{}{
			"servers": []map[string]interface{}{
				{
					"tag":     "remote",
					"address": "https://1.1.1.1/dns-query",
					"detour":  "proxy",
				},
				{
					"tag":     "local",
					"address": "https://223.5.5.5/dns-query",
					"detour":  "direct",
				},
			},
			"rules": []map[string]interface{}{
				{
					"geosite": []string{"cn"},
					"server":  "local",
				},
				{
					"geosite": []string{"geolocation-!cn"},
					"server":  "remote",
				},
			},
			"final": "remote",
		},
		"inbounds": []map[string]interface{}{
			{
				"type":        "mixed",
				"tag":         "mixed-in",
				"listen":      "127.0.0.1",
				"listen_port": 7890,
				"sniff":       true,
			},
			{
				"type":                    "tun",
				"tag":                     "tun-in",
				"inet4_address":           "172.19.0.1/30",
				"auto_route":              true,
				"strict_route":            true,
				"sniff":                   true,
				"sniff_override_destination": true,
			},
		},
		"outbounds": []map[string]interface{}{},
		"route": map[string]interface{}{
			"rules": []map[string]interface{}{
				{
					"protocol": "dns",
					"outbound": "dns-out",
				},
				{
					"geosite":  []string{"cn", "private"},
					"geoip":    []string{"cn", "private"},
					"outbound": "direct",
				},
				{
					"geosite":  []string{"geolocation-!cn"},
					"outbound": "proxy",
				},
			},
			"final":                 "proxy",
			"auto_detect_interface": true,
		},
	}

	// 添加出站节点
	outbounds := config["outbounds"].([]map[string]interface{})

	// 添加必要的出站
	outbounds = append(outbounds,
		map[string]interface{}{
			"type": "direct",
			"tag":  "direct",
		},
		map[string]interface{}{
			"type": "block",
			"tag":  "block",
		},
		map[string]interface{}{
			"type": "dns",
			"tag":  "dns-out",
		},
	)

	// 转换服务器节点
	var proxyTags []string
	for _, server := range sub.Servers {
		node := server.ConvertToSingboxNode()
		outbounds = append(outbounds, node)
		proxyTags = append(proxyTags, server.Name)
	}

	// 添加选择器
	if len(proxyTags) > 0 {
		// 自动选择
		outbounds = append(outbounds, map[string]interface{}{
			"type":      "urltest",
			"tag":       "auto",
			"outbounds": proxyTags,
			"url":       "http://www.gstatic.com/generate_204",
			"interval":  "5m",
			"tolerance": 50,
		})

		// 手动选择
		outbounds = append(outbounds, map[string]interface{}{
			"type":      "selector",
			"tag":       "select",
			"outbounds": append([]string{"auto"}, proxyTags...),
			"default":   "auto",
		})

		// 主代理选择器
		outbounds = append(outbounds, map[string]interface{}{
			"type":      "selector",
			"tag":       "proxy",
			"outbounds": []string{"select", "direct"},
			"default":   "select",
		})
	}

	config["outbounds"] = outbounds

	return config, nil
}

// ParseSubscriptionURL 解析订阅 URL，提取 baseURL 和 token
func ParseSubscriptionURL(url string) (baseURL, token string, err error) {
	// xboard 订阅 URL 格式通常为：https://example.com/api/v1/client/subscribe?token=xxx
	// 或者：https://example.com/sub/xxx

	if strings.Contains(url, "/api/v1/client/subscribe") {
		// 标准 API 格式
		parts := strings.Split(url, "?")
		if len(parts) < 2 {
			return "", "", fmt.Errorf("无效的订阅 URL")
		}

		baseURL = strings.Replace(parts[0], "/api/v1/client/subscribe", "", 1)

		// 解析 token
		params := strings.Split(parts[1], "&")
		for _, param := range params {
			if strings.HasPrefix(param, "token=") {
				token = strings.TrimPrefix(param, "token=")
				break
			}
		}
	} else if strings.Contains(url, "/sub/") {
		// 短链接格式
		parts := strings.Split(url, "/sub/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("无效的订阅 URL")
		}

		baseURL = parts[0]
		token = parts[1]
	} else {
		return "", "", fmt.Errorf("不支持的订阅 URL 格式")
	}

	if baseURL == "" || token == "" {
		return "", "", fmt.Errorf("无法解析订阅 URL")
	}

	return baseURL, token, nil
}

// SetLogger 设置日志记录器
func (c *Client) SetLogger(logger *logrus.Logger) {
	c.logger = logger
}