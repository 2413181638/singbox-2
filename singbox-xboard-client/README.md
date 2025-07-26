# Singbox Xboard Client

基于 sing-box 内核的跨平台客户端，对接 xboard 面板。

## 功能特性

- ✅ 支持 xboard 面板订阅
- ✅ 支持最新协议：Hysteria2、VLESS+Reality、Shadowsocks 2022
- ✅ 跨平台支持：Windows、macOS、Android
- ✅ 自动更新订阅
- ✅ 流量统计
- ✅ 规则分流
- ✅ GitHub Actions 自动打包

## 支持的协议

- VLESS + Reality + XTLS
- Hysteria2
- Shadowsocks 2022
- VMess
- Trojan
- TUIC v5

## 快速开始

### 下载安装

从 [Releases](https://github.com/your-username/singbox-xboard-client/releases) 页面下载对应平台的安装包。

### 配置使用

1. 启动客户端
2. 输入 xboard 订阅地址
3. 选择节点连接

## 开发说明

### 环境要求

- Go 1.21+
- Node.js 18+
- Android Studio (Android 开发)
- Xcode (macOS/iOS 开发)

### 构建项目

```bash
# 克隆项目
git clone https://github.com/your-username/singbox-xboard-client.git
cd singbox-xboard-client

# 安装依赖
make deps

# 构建所有平台
make all

# 构建特定平台
make windows
make macos
make android
```

## 项目结构

```
singbox-xboard-client/
├── cmd/                    # 命令行入口
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── subscription/     # 订阅管理
│   ├── singbox/         # sing-box 集成
│   └── ui/              # 用户界面
├── pkg/                   # 公共包
├── build/                # 构建脚本
├── assets/              # 资源文件
└── .github/             # GitHub Actions
```

## 开源协议

GPL-3.0 License