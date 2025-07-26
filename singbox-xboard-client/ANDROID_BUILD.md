# Android APK 构建和发布指南

## 自动构建流程

当你将代码推送到 GitHub 并创建标签时，GitHub Actions 会自动构建并发布 Android APK。

### 1. 推送代码到 GitHub

```bash
cd /workspace/singbox-xboard-client
git init
git add .
git commit -m "Initial commit with Android support"
git remote add origin https://github.com/YOUR_USERNAME/singbox-xboard-client.git
git push -u origin main
```

### 2. 创建发布标签

```bash
git tag v1.0.0
git push origin v1.0.0
```

### 3. 自动构建和发布

GitHub Actions 会自动：
- 设置 Android SDK 环境
- 安装必要的构建工具
- 使用 gomobile 构建 sing-box 核心库
- 构建 Android APK
- 签名 APK（使用 CI 密钥）
- 发布到 GitHub Releases

### 4. 下载 APK

构建完成后（约 10-15 分钟），你可以：
1. 访问你的 GitHub 仓库
2. 点击 "Releases" 标签
3. 找到最新版本
4. 下载 `singbox-xboard-android.apk`

## 本地构建（可选）

如果你想在本地构建 APK：

### 环境要求

1. **Android SDK**
   ```bash
   # 安装 Android Studio 或 Android Command Line Tools
   export ANDROID_HOME=/path/to/android-sdk
   export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools
   ```

2. **Go 1.21+**
   ```bash
   go version  # 应该显示 go1.21 或更高
   ```

3. **Java 17**
   ```bash
   java -version  # 应该显示 17
   ```

### 构建步骤

```bash
cd /workspace/singbox-xboard-client

# 安装依赖
make deps

# 构建 Android APK
make android

# APK 文件位置
ls build/singbox-xboard-android.apk
```

## APK 功能

构建的 APK 包含：

- ✅ 完整的 sing-box 核心功能
- ✅ xboard 面板订阅支持
- ✅ 支持所有协议（Hysteria2、VLESS+Reality 等）
- ✅ VPN 服务实现
- ✅ Material Design UI
- ✅ 自动订阅更新
- ✅ 流量统计
- ✅ 节点切换

## 安装 APK

1. 在 Android 设备上下载 APK
2. 打开文件管理器，找到下载的 APK
3. 点击 APK 文件
4. 如果提示，允许"安装未知来源应用"
5. 按照提示完成安装
6. 首次运行时授予 VPN 权限

## 使用说明

1. **添加订阅**
   - 打开应用
   - 输入 xboard 订阅地址
   - 点击"更新订阅"

2. **连接 VPN**
   - 选择节点（可选）
   - 点击"连接"按钮
   - 授予 VPN 权限（首次）

3. **查看状态**
   - 主界面显示连接状态
   - 实时流量统计
   - 用户信息和流量使用情况

## 故障排除

### 构建失败

1. **ANDROID_HOME 未设置**
   ```bash
   export ANDROID_HOME=/path/to/android-sdk
   ```

2. **gomobile 初始化失败**
   ```bash
   go clean -cache
   go install golang.org/x/mobile/cmd/gomobile@latest
   gomobile init
   ```

3. **Gradle 构建失败**
   ```bash
   cd android
   ./gradlew clean
   ./gradlew assembleRelease
   ```

### 运行问题

1. **无法安装 APK**
   - 确保允许"安装未知来源应用"
   - 检查 Android 版本（需要 5.0+）

2. **VPN 连接失败**
   - 检查订阅地址是否正确
   - 确保网络连接正常
   - 查看应用日志

## 开发说明

### 项目结构

```
android/
├── app/
│   ├── src/main/java/com/singbox/xboard/
│   │   ├── MainActivity.java          # 主界面
│   │   ├── service/
│   │   │   └── VpnService.java       # VPN 服务
│   │   └── viewmodel/
│   │       └── MainViewModel.java     # 数据管理
│   └── build.gradle                   # 应用配置
├── build.gradle                       # 项目配置
└── gradlew                           # 构建脚本
```

### 核心组件

1. **pkg/mobile/mobile.go**
   - Go 代码接口，供 Android 调用
   - 管理 sing-box 核心功能

2. **VpnService**
   - Android VPN 服务实现
   - 处理网络流量转发

3. **MainActivity**
   - 用户界面
   - 订阅管理
   - 连接控制

## 签名说明

GitHub Actions 构建的 APK 使用临时签名密钥。如果需要发布到 Google Play：

1. 生成正式签名密钥
2. 在 GitHub Secrets 中配置签名信息
3. 更新工作流使用正式签名

## 更新日志

### v1.0.0
- 初始版本
- 支持 xboard 订阅
- 实现基本 VPN 功能
- Material Design UI