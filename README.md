# SingBox App

基于 sing-box 的代理应用程序，支持多种协议和配置方式，提供CLI和Web两种管理界面。

## 功能特性

- 🚀 基于 sing-box 核心，性能优异
- 🎯 支持多种代理协议（SOCKS5、HTTP等）
- 🖥️ 提供CLI和Web两种管理界面
- ⚙️ 灵活的配置管理
- 📊 实时状态监控
- 🎨 现代化Web界面设计

## 快速开始

### 环境要求

- Go 1.21+
- Linux/macOS/Windows

### 编译安装

```bash
# 克隆项目
git clone <repository-url>
cd singbox-app

# 安装依赖
make deps

# 编译
make build

# 运行
./build/singbox-app
```

### 使用方法

#### CLI模式

```bash
# 使用默认配置启动
./build/singbox-app

# 指定配置文件
./build/singbox-app -c /path/to/config.yaml

# 查看帮助
./build/singbox-app --help
```

#### Web界面模式

```bash
# 启动Web界面
./build/singbox-app --web --port 8080

# 访问管理界面
# 浏览器打开: http://localhost:8080
```

## 配置说明

应用程序会自动创建默认配置文件 `config.yaml`：

```yaml
inbound:
  type: socks          # 入站类型: socks, http
  port: 1080          # 监听端口
  host: 127.0.0.1     # 监听地址

outbound:
  type: direct        # 出站类型: direct, socks, shadowsocks等
  server: ""          # 服务器地址
  port: 0             # 服务器端口

dns:
  servers:
    - 8.8.8.8
    - 1.1.1.1

log:
  level: info         # 日志级别: debug, info, warn, error
```

## 开发指南

### 项目结构

```
singbox-app/
├── cmd/                 # 命令行相关
│   ├── root.go         # 根命令
│   ├── cli.go          # CLI模式
│   └── web.go          # Web模式
├── internal/           # 内部包
│   ├── config/         # 配置管理
│   ├── logger/         # 日志管理
│   ├── proxy/          # 代理服务
│   └── web/            # Web界面
├── web/                # Web资源
│   └── templates/      # HTML模板
├── main.go             # 主入口
├── go.mod              # Go模块
├── Makefile            # 构建脚本
└── README.md           # 说明文档
```

### 开发命令

```bash
# 开发模式（热重载）
make dev

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make vet

# 交叉编译
make build-all
```

### 添加新功能

1. 在相应的模块中添加功能代码
2. 更新配置结构（如需要）
3. 添加Web API接口（如需要）
4. 更新前端界面（如需要）
5. 编写测试用例
6. 更新文档

## API文档

### REST API

- `GET /api/status` - 获取服务状态
- `GET /api/config` - 获取当前配置
- `POST /api/config` - 更新配置
- `POST /api/start` - 启动服务
- `POST /api/stop` - 停止服务

### 响应格式

```json
{
  "status": "ok",
  "running": true,
  "stats": {
    "connections": 0,
    "traffic": {
      "upload": 0,
      "download": 0
    }
  }
}
```

## 部署说明

### 系统服务部署

```bash
# 复制可执行文件
sudo cp build/singbox-app /usr/local/bin/

# 创建配置目录
sudo mkdir -p /etc/singbox-app

# 复制配置文件
sudo cp config.yaml /etc/singbox-app/

# 创建systemd服务文件
sudo tee /etc/systemd/system/singbox-app.service > /dev/null <<EOF
[Unit]
Description=SingBox App
After=network.target

[Service]
Type=simple
User=nobody
ExecStart=/usr/local/bin/singbox-app -c /etc/singbox-app/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动服务
sudo systemctl enable singbox-app
sudo systemctl start singbox-app
```

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o singbox-app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/singbox-app .
COPY --from=builder /app/config.yaml .
EXPOSE 1080 8080
CMD ["./singbox-app", "--web"]
```

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 检查端口占用
   netstat -tlnp | grep :1080
   
   # 修改配置文件中的端口
   vim config.yaml
   ```

2. **权限不足**
   ```bash
   # 确保有执行权限
   chmod +x build/singbox-app
   
   # 使用非特权端口（>1024）
   ```

3. **配置文件错误**
   ```bash
   # 验证YAML格式
   ./build/singbox-app -c config.yaml --dry-run
   ```

### 日志查看

```bash
# 实时查看日志
tail -f /var/log/singbox-app.log

# 查看systemd日志
sudo journalctl -u singbox-app -f
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 支持

如果您遇到问题或有建议，请：

1. 查看 [Issues](../../issues) 页面
2. 创建新的 Issue
3. 联系维护者

---

**注意**: 本应用程序仅供学习和研究使用，请遵守当地法律法规。
