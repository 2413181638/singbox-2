# 部署指南

## 快速开始

### 1. 克隆仓库

```bash
git clone https://github.com/your-username/singbox-xboard-client.git
cd singbox-xboard-client
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 构建项目

#### 构建当前平台
```bash
make build
```

#### 构建所有平台
```bash
make release
```

#### 构建特定平台
```bash
# Windows
make windows

# macOS
make macos

# Linux
make linux
```

### 4. 运行

```bash
./build/singbox-xboard
```

## GitHub Actions 自动构建

本项目已配置 GitHub Actions，当你：

1. Push 到 `main` 分支时，会自动构建
2. 创建标签（如 `v1.0.0`）时，会自动创建 Release 并上传构建产物

### 创建 Release

```bash
# 创建标签
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 会自动：
- 构建所有平台的二进制文件
- 创建 Release
- 上传构建产物
- 生成 SHA256 校验和

## 手动部署

### Windows

1. 下载 `singbox-xboard-windows-amd64.exe`
2. 双击运行
3. 浏览器会自动打开 `http://localhost:9090`

### macOS

1. 下载对应版本：
   - Intel: `singbox-xboard-darwin-amd64`
   - Apple Silicon: `singbox-xboard-darwin-arm64`
2. 赋予执行权限：`chmod +x singbox-xboard-darwin-*`
3. 运行：`./singbox-xboard-darwin-*`

### Linux

1. 下载对应架构的版本
2. 赋予执行权限：`chmod +x singbox-xboard-linux-*`
3. 运行：`./singbox-xboard-linux-*`

### Docker（计划中）

```bash
docker run -d \
  --name singbox-xboard \
  -p 9090:9090 \
  -p 7890:7890 \
  -v /path/to/config:/config \
  ghcr.io/your-username/singbox-xboard:latest
```

## 配置说明

首次运行时会自动创建默认配置文件，位置：
- Windows: `%APPDATA%\singbox-xboard\config.json`
- macOS: `~/Library/Application Support/singbox-xboard/config.json`
- Linux: `~/.config/singbox-xboard/config.json`

## 故障排除

### 端口被占用

修改配置文件中的端口设置：

```json
{
  "ui": {
    "port": 9091
  },
  "singbox": {
    "inbounds": [{
      "listen_port": 7891
    }]
  }
}
```

### sing-box 核心找不到

运行下载脚本：
```bash
make download-singbox
```

或手动下载 sing-box 并放置到：
- 程序同目录
- `bin/` 子目录
- 系统 PATH 中

## 更多信息

- [使用文档](README.md)
- [更新日志](CHANGELOG.md)
- [问题反馈](https://github.com/your-username/singbox-xboard-client/issues)