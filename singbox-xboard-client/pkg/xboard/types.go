package xboard

import (
	"fmt"
	"time"
)

// SubscriptionResponse xboard 订阅响应
type SubscriptionResponse struct {
	Servers []Server `json:"servers"`
}

// Server 服务器节点信息
type Server struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	Type       string   `json:"type"`       // shadowsocks, vmess, vless, trojan, hysteria2
	Cipher     string   `json:"cipher"`     // 加密方式
	UUID       string   `json:"uuid"`       // UUID
	Password   string   `json:"password"`   // 密码
	AlterId    int      `json:"alter_id"`   // VMess alterID
	Network    string   `json:"network"`    // 传输协议: tcp, ws, grpc, quic
	Path       string   `json:"path"`       // WebSocket/gRPC 路径
	TLS        bool     `json:"tls"`        // 是否启用 TLS
	SkipCert   bool     `json:"skip_cert"`  // 跳过证书验证
	ServerName string   `json:"sni"`        // SNI
	Flow       string   `json:"flow"`       // VLESS flow
	Tags       []string `json:"tags"`       // 标签
	
	// Reality 配置
	Reality *RealityConfig `json:"reality,omitempty"`
	
	// Hysteria2 配置
	Hysteria2 *Hysteria2Config `json:"hysteria2,omitempty"`
}

// RealityConfig Reality 配置
type RealityConfig struct {
	PublicKey  string `json:"public_key"`  // 公钥
	ShortID    string `json:"short_id"`    // 短 ID
	ServerName string `json:"server_name"` // 目标服务器名称
}

// Hysteria2Config Hysteria2 配置
type Hysteria2Config struct {
	Up   string `json:"up"`   // 上行带宽
	Down string `json:"down"` // 下行带宽
	Obfs string `json:"obfs"` // 混淆密码
}

// UserInfo 用户信息
type UserInfo struct {
	Email        string    `json:"email"`
	Upload       int64     `json:"upload"`       // 已上传流量（字节）
	Download     int64     `json:"download"`     // 已下载流量（字节）
	Total        int64     `json:"total"`        // 总流量（字节）
	ExpireTime   time.Time `json:"expire_time"`  // 过期时间
	DeviceLimit  int       `json:"device_limit"` // 设备限制
	SpeedLimit   int       `json:"speed_limit"`  // 速度限制（Mbps）
}

// NodeInfo 节点信息
type NodeInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Info     string `json:"info"`
	Status   int    `json:"status"`   // 0: 离线, 1: 在线
	Load     int    `json:"load"`     // 负载百分比
	Uptime   int64  `json:"uptime"`   // 运行时间（秒）
	Network  string `json:"network"`  // 网络类型
	Location string `json:"location"` // 位置
}

// TrafficLog 流量日志
type TrafficLog struct {
	UserID     int       `json:"user_id"`
	Upload     int64     `json:"upload"`
	Download   int64     `json:"download"`
	NodeID     int       `json:"node_id"`
	Rate       float64   `json:"rate"`
	LogTime    time.Time `json:"log_time"`
}

// APIError API 错误
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ConvertToSingboxNode 转换为 sing-box 节点配置
func (s *Server) ConvertToSingboxNode() map[string]interface{} {
	node := make(map[string]interface{})
	node["tag"] = s.Name
	node["type"] = s.Type
	
	switch s.Type {
	case "shadowsocks":
		node["method"] = s.Cipher
		node["password"] = s.Password
		node["server"] = s.Host
		node["server_port"] = s.Port
		
	case "vmess":
		node["server"] = s.Host
		node["server_port"] = s.Port
		node["uuid"] = s.UUID
		node["alter_id"] = s.AlterId
		node["security"] = s.Cipher
		
		if s.Network != "tcp" {
			transport := make(map[string]interface{})
			transport["type"] = s.Network
			
			if s.Network == "ws" {
				transport["path"] = s.Path
				if s.Host != "" {
					headers := make(map[string]string)
					headers["Host"] = s.Host
					transport["headers"] = headers
				}
			} else if s.Network == "grpc" {
				transport["service_name"] = s.Path
			}
			
			node["transport"] = transport
		}
		
		if s.TLS {
			tls := make(map[string]interface{})
			tls["enabled"] = true
			if s.ServerName != "" {
				tls["server_name"] = s.ServerName
			}
			tls["insecure"] = s.SkipCert
			node["tls"] = tls
		}
		
	case "vless":
		node["server"] = s.Host
		node["server_port"] = s.Port
		node["uuid"] = s.UUID
		node["flow"] = s.Flow
		
		if s.Reality != nil {
			tls := make(map[string]interface{})
			tls["enabled"] = true
			tls["server_name"] = s.Reality.ServerName
			
			reality := make(map[string]interface{})
			reality["enabled"] = true
			reality["public_key"] = s.Reality.PublicKey
			reality["short_id"] = s.Reality.ShortID
			
			tls["reality"] = reality
			node["tls"] = tls
		} else if s.TLS {
			tls := make(map[string]interface{})
			tls["enabled"] = true
			if s.ServerName != "" {
				tls["server_name"] = s.ServerName
			}
			tls["insecure"] = s.SkipCert
			node["tls"] = tls
		}
		
		if s.Network != "tcp" {
			transport := make(map[string]interface{})
			transport["type"] = s.Network
			
			if s.Network == "ws" {
				transport["path"] = s.Path
				if s.Host != "" {
					headers := make(map[string]string)
					headers["Host"] = s.Host
					transport["headers"] = headers
				}
			} else if s.Network == "grpc" {
				transport["service_name"] = s.Path
			}
			
			node["transport"] = transport
		}
		
	case "trojan":
		node["server"] = s.Host
		node["server_port"] = s.Port
		node["password"] = s.Password
		
		tls := make(map[string]interface{})
		tls["enabled"] = true
		if s.ServerName != "" {
			tls["server_name"] = s.ServerName
		}
		tls["insecure"] = s.SkipCert
		node["tls"] = tls
		
		if s.Network != "tcp" {
			transport := make(map[string]interface{})
			transport["type"] = s.Network
			
			if s.Network == "ws" {
				transport["path"] = s.Path
				if s.Host != "" {
					headers := make(map[string]string)
					headers["Host"] = s.Host
					transport["headers"] = headers
				}
			} else if s.Network == "grpc" {
				transport["service_name"] = s.Path
			}
			
			node["transport"] = transport
		}
		
	case "hysteria2":
		node["server"] = s.Host
		node["server_port"] = s.Port
		node["password"] = s.Password
		
		if s.Hysteria2 != nil {
			if s.Hysteria2.Up != "" {
				node["up_mbps"] = parseSpeed(s.Hysteria2.Up)
			}
			if s.Hysteria2.Down != "" {
				node["down_mbps"] = parseSpeed(s.Hysteria2.Down)
			}
			if s.Hysteria2.Obfs != "" {
				obfs := make(map[string]interface{})
				obfs["type"] = "salamander"
				obfs["password"] = s.Hysteria2.Obfs
				node["obfs"] = obfs
			}
		}
		
		tls := make(map[string]interface{})
		tls["enabled"] = true
		if s.ServerName != "" {
			tls["server_name"] = s.ServerName
		}
		tls["insecure"] = s.SkipCert
		node["tls"] = tls
	}
	
	return node
}

// parseSpeed 解析速度字符串（如 "100 Mbps" -> 100）
func parseSpeed(speed string) int {
	// 简单实现，实际应该更复杂
	var value int
	fmt.Sscanf(speed, "%d", &value)
	return value
}