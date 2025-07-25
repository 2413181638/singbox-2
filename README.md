# SingBox XBoard Client

一个基于 SingBox 内核的现代化代理客户端，完美对接 XBoard 面板，提供优雅的用户界面和强大的功能。

## 特性

- 🚀 **基于 SingBox 内核** - 高性能、低延迟的代理核心
- 🎯 **完美对接 XBoard** - 无缝集成 XBoard 面板，自动同步订阅
- 🖥️ **现代化界面** - 基于 Vue 3 + Element Plus 的美观界面
- 🔄 **自动订阅同步** - 定时同步服务器节点，保持最新状态
- 📊 **实时流量统计** - 详细的流量使用情况和连接状态监控
- 🌐 **多平台支持** - Windows、macOS、Linux 全平台支持
- 🐳 **Docker 支持** - 支持容器化部署
- 🔧 **灵活配置** - 丰富的配置选项，满足不同需求

## 支持的协议

- VMess
- VLESS
- Shadowsocks
- Trojan
- Hysteria

## 快速开始

### 下载安装

从 [Releases](https://github.com/your-username/singbox-xboard-client/releases) 页面下载对应平台的客户端：

- **Windows**: `singbox-xboard-client-windows-amd64.zip`
- **macOS**: `singbox-xboard-client-darwin-amd64.tar.gz` (Intel) / `singbox-xboard-client-darwin-arm64.tar.gz` (Apple Silicon)
- **Linux**: `singbox-xboard-client-linux-amd64.tar.gz`

### 使用方法

1. **启动应用**
   ```bash
   # 解压后直接运行
   ./singbox-xboard-client
   ```

2. **登录 XBoard 面板**
   - 输入 XBoard 面板地址
   - 输入邮箱和密码
   - 点击登录

3. **开始使用**
   - 应用会自动同步服务器节点
   - 在仪表板查看连接状态和流量统计
   - 在服务器页面管理节点
   - 点击连接按钮开始代理

### 代理设置

应用启动后，会在本地创建以下代理端口：

- **HTTP 代理**: `127.0.0.1:7890`
- **SOCKS5 代理**: `127.0.0.1:7890`

将这些地址配置到您的应用程序中即可使用代理。

## 开发

### 环境要求

- Go 1.21+
- Node.js 18+
- Wails v2

### 安装依赖

```bash
# 安装 Go 依赖
go mod download

# 安装前端依赖
cd frontend
npm install
```

### 开发运行

```bash
# 开发模式运行
wails dev
```

### 构建

```bash
# 构建当前平台
wails build

# 构建指定平台
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64
```

## Docker 部署

### 使用 Docker Hub 镜像

```bash
docker run -d \
  --name singbox-xboard-client \
  -p 7890:7890 \
  -p 9090:9090 \
  -v /path/to/config:/app/config \
  singbox-xboard-client:latest
```

### 从源码构建

```bash
# 构建镜像
docker build -t singbox-xboard-client .

# 运行容器
docker run -d \
  --name singbox-xboard-client \
  -p 7890:7890 \
  -p 9090:9090 \
  singbox-xboard-client
```

## 配置文件

配置文件位于 `~/.singbox-xboard/config.yaml`：

```yaml
database_path: ~/.singbox-xboard/data.db
log_level: info
xboard:
  url: "https://your-xboard-panel.com"
  token: "your-token"
  node_id: 1
  interval: 300
singbox:
  config_path: ~/.singbox-xboard/singbox.json
  log_path: ~/.singbox-xboard/singbox.log
  api_port: 9090
```

## 功能截图

### 仪表板
- 用户信息展示
- 流量使用统计
- 连接状态监控
- 服务器状态概览

### 服务器管理
- 服务器列表展示
- 延迟测试
- 批量操作
- 服务器详情查看

### 连接日志
- 详细的连接记录
- 流量统计
- 错误信息追踪

### 设置页面
- 基本配置管理
- 系统信息查看
- 配置导入导出

## API 接口

应用提供 RESTful API 接口，方便第三方集成：

- `GET /api/status` - 获取连接状态
- `POST /api/start` - 启动连接
- `POST /api/stop` - 停止连接
- `GET /api/servers` - 获取服务器列表
- `POST /api/servers/sync` - 同步订阅

## 故障排除

### 常见问题

1. **无法连接到 XBoard 面板**
   - 检查面板地址是否正确
   - 确认网络连接正常
   - 检查防火墙设置

2. **代理无法使用**
   - 确认应用已启动连接
   - 检查代理端口是否被占用
   - 查看日志文件排查错误

3. **服务器延迟测试失败**
   - 检查服务器配置是否正确
   - 确认服务器可达性
   - 查看错误日志

### 日志文件

- 应用日志: `~/.singbox-xboard/app.log`
- SingBox 日志: `~/.singbox-xboard/singbox.log`

## 贡献

欢迎提交 Issue 和 Pull Request！

### 开发指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 致谢

- [SingBox](https://github.com/SagerNet/sing-box) - 强大的代理核心
- [Wails](https://github.com/wailsapp/wails) - 优秀的 Go + Web 桌面应用框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Element Plus](https://element-plus.org/) - 基于 Vue 3 的组件库

## 支持

如果这个项目对您有帮助，请考虑给它一个 ⭐️！

## 联系方式

- 项目主页: https://github.com/your-username/singbox-xboard-client
- 问题反馈: https://github.com/your-username/singbox-xboard-client/issues
- 讨论交流: https://github.com/your-username/singbox-xboard-client/discussions
