# SingBox App 使用指南

这是一个基于 sing-box 的代理应用程序，提供了完整的开发框架和现代化的管理界面。

## 🚀 快速开始

### 1. 编译应用程序

```bash
# 安装依赖
make deps

# 编译
make build

# 或者一键编译
make all
```

### 2. 运行应用程序

#### CLI模式
```bash
# 使用默认配置
./build/singbox-app

# 指定配置文件
./build/singbox-app -c examples/config-socks.yaml

# 查看帮助
./build/singbox-app --help
```

#### Web界面模式
```bash
# 启动Web管理界面
./build/singbox-app --web --port 8080

# 访问: http://localhost:8080
```

### 3. 使用测试脚本

```bash
# CLI模式测试
./scripts/test-cli.sh

# Web界面测试
./scripts/test-web.sh

# 多平台编译
./scripts/build-all.sh
```

## 📋 配置说明

### 基本配置结构

```yaml
inbound:
  type: socks          # 入站类型: socks, http
  port: 1080          # 监听端口
  host: 127.0.0.1     # 监听地址

outbound:
  type: direct        # 出站类型: direct, socks
  server: ""          # 服务器地址
  port: 0             # 服务器端口
  username: ""        # 用户名（可选）
  password: ""        # 密码（可选）

dns:
  servers:
    - 8.8.8.8         # DNS服务器列表
    - 1.1.1.1

log:
  level: info         # 日志级别: debug, info, warn, error
  file: ""            # 日志文件（可选）
```

### 示例配置

#### SOCKS5 本地代理
```bash
cp examples/config-socks.yaml config.yaml
./build/singbox-app -c config.yaml
```

#### HTTP 本地代理
```bash
cp examples/config-http.yaml config.yaml
./build/singbox-app -c config.yaml
```

#### SOCKS5 上游代理
```bash
cp examples/config-socks-upstream.yaml config.yaml
# 编辑配置文件，设置上游服务器信息
./build/singbox-app -c config.yaml
```

## 🌐 Web管理界面

### 功能特性

- 📊 实时服务状态监控
- ⚙️ 在线配置管理
- 🎛️ 服务启停控制
- 📝 实时日志显示
- 🎨 现代化响应式界面

### 使用说明

1. **启动Web界面**
   ```bash
   ./build/singbox-app --web --port 8080
   ```

2. **访问管理面板**
   - 打开浏览器访问: http://localhost:8080
   - 界面包含服务状态、配置管理、日志显示等功能

3. **API接口**
   - `GET /api/status` - 获取服务状态
   - `POST /api/start` - 启动服务
   - `POST /api/stop` - 停止服务
   - `GET /api/config` - 获取配置
   - `POST /api/config` - 更新配置

## 🐳 Docker 部署

### 快速部署

```bash
# 一键Docker部署
./scripts/docker-run.sh
```

### 手动部署

```bash
# 构建镜像
docker build -t singbox-app .

# 运行容器
docker run -d \
  --name singbox-app \
  -p 1080:1080 \
  -p 8080:8080 \
  -v $(pwd)/config:/etc/singbox-app \
  singbox-app
```

### Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 🔧 开发指南

### 项目结构

```
singbox-app/
├── cmd/                 # 命令行接口
├── internal/           # 内部包
│   ├── config/         # 配置管理
│   ├── logger/         # 日志系统
│   ├── proxy/          # 代理服务
│   └── web/            # Web界面
├── web/                # Web资源
├── scripts/            # 脚本工具
├── examples/           # 示例配置
├── build/              # 编译输出
└── docs/               # 文档
```

### 开发命令

```bash
# 开发模式（热重载）
make dev

# 代码格式化
make fmt

# 代码检查
make vet

# 运行测试
make test

# 清理构建文件
make clean
```

### 添加新功能

1. **添加新的入站/出站类型**
   - 修改 `internal/config/config.go`
   - 更新 `internal/proxy/service.go`

2. **扩展Web API**
   - 修改 `internal/web/routes.go`
   - 更新前端界面 `web/templates/index.html`

3. **集成真正的sing-box**
   - 参考项目中的注释和TODO
   - 使用sing-box的官方API

## 📖 常见问题

### Q: 如何更改监听端口？
A: 修改配置文件中的 `inbound.port` 字段，或使用 `--port` 参数。

### Q: 如何添加认证？
A: 当前版本为演示版本，认证功能需要集成完整的sing-box实现。

### Q: 如何查看详细日志？
A: 设置配置文件中的 `log.level` 为 `debug`，或查看Web界面的日志面板。

### Q: 如何部署到生产环境？
A: 建议使用Docker部署，或者编译后使用systemd管理服务。

## 🔒 安全注意事项

1. **网络安全**
   - 仅在受信任的网络环境中使用
   - 生产环境请配置防火墙规则

2. **配置安全**
   - 保护配置文件中的敏感信息
   - 使用强密码和加密连接

3. **访问控制**
   - Web界面建议配置访问限制
   - 考虑添加身份验证

## 📞 支持与反馈

- 📧 问题报告: 创建GitHub Issue
- 💡 功能建议: 提交Pull Request
- 📚 文档改进: 欢迎贡献

## 📄 许可证

本项目采用 MIT 许可证，详见 LICENSE 文件。

---

**注意**: 本应用程序仅供学习和研究使用，请遵守当地法律法规。