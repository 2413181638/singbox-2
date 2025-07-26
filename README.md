# Singbox XBoard 客户端

基于 sing-box 内核的跨平台 xboard 面板客户端，支持一键获取配置、启动/停止 sing-box。

## 功能
- 对接 xboard 面板，自动获取节点配置
- 一键启动/停止 sing-box
- 跨平台（Windows/macOS/Linux）
- GitHub Actions 自动打包发布

## 使用方法

1. 克隆本仓库，安装依赖：
   ```bash
   npm install
   ```
2. 构建前端：
   ```bash
   npm run build
   ```
3. 启动客户端：
   ```bash
   npm start
   ```
4. 输入 xboard token，获取配置并启动 sing-box

## 打包与发布

- 推送 tag（如 v0.1.0）到 GitHub，Actions 会自动打包并发布 release。

## 目录结构
- `main/`：Electron 主进程、sing-box、xboard 对接
- `src/`：React 前端
- `sing-box/`：存放 sing-box 各平台二进制

## 备注
- 请自行下载 sing-box 各平台二进制，放入 `sing-box/` 目录。
- xboard API 地址和 token 获取方式请参考你的面板文档。
