# 变量定义
BINARY_NAME=singbox-xboard
VERSION?=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Go 相关
GO=go
GOFLAGS=-trimpath
GOBUILD=$(GO) build $(GOFLAGS)
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get
GOMOD=$(GO) mod

# 平台定义
PLATFORMS=windows/amd64 windows/386 darwin/amd64 darwin/arm64 linux/amd64 linux/386 linux/arm linux/arm64

# 输出目录
DIST_DIR=dist
BUILD_DIR=build

.PHONY: all build clean test deps help

# 默认目标
all: deps build

# 显示帮助
help:
	@echo "Makefile for $(BINARY_NAME)"
	@echo ""
	@echo "Usage:"
	@echo "  make deps       安装依赖"
	@echo "  make build      构建当前平台"
	@echo "  make all        安装依赖并构建"
	@echo "  make windows    构建 Windows 版本"
	@echo "  make macos      构建 macOS 版本"
	@echo "  make linux      构建 Linux 版本"
	@echo "  make android    构建 Android 版本"
	@echo "  make release    构建所有平台发布版本"
	@echo "  make clean      清理构建文件"
	@echo "  make test       运行测试"

# 安装依赖
deps:
	@echo "==> 安装依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 构建当前平台
build:
	@echo "==> 构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

# 构建 Windows
windows:
	@echo "==> 构建 Windows 版本..."
	@mkdir -p $(DIST_DIR)/windows-amd64
	@mkdir -p $(DIST_DIR)/windows-386
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/windows-amd64/$(BINARY_NAME).exe ./cmd
	GOOS=windows GOARCH=386 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/windows-386/$(BINARY_NAME).exe ./cmd
	@echo "Windows 版本构建完成"

# 构建 macOS
macos:
	@echo "==> 构建 macOS 版本..."
	@mkdir -p $(DIST_DIR)/darwin-amd64
	@mkdir -p $(DIST_DIR)/darwin-arm64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/darwin-amd64/$(BINARY_NAME) ./cmd
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/darwin-arm64/$(BINARY_NAME) ./cmd
	@echo "macOS 版本构建完成"

# 构建 Linux
linux:
	@echo "==> 构建 Linux 版本..."
	@mkdir -p $(DIST_DIR)/linux-amd64
	@mkdir -p $(DIST_DIR)/linux-386
	@mkdir -p $(DIST_DIR)/linux-arm
	@mkdir -p $(DIST_DIR)/linux-arm64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/linux-amd64/$(BINARY_NAME) ./cmd
	GOOS=linux GOARCH=386 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/linux-386/$(BINARY_NAME) ./cmd
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/linux-arm/$(BINARY_NAME) ./cmd
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/linux-arm64/$(BINARY_NAME) ./cmd
	@echo "Linux 版本构建完成"

# 构建 Android (需要 Android NDK)
android:
	@echo "==> 构建 Android 版本..."
	@if [ -z "$(ANDROID_HOME)" ]; then \
		echo "错误: ANDROID_HOME 环境变量未设置"; \
		exit 1; \
	fi
	@echo "安装 gomobile..."
	go install golang.org/x/mobile/cmd/gomobile@latest
	go install golang.org/x/mobile/cmd/gobind@latest
	@echo "初始化 gomobile..."
	gomobile init || true
	@echo "构建 Android AAR..."
	cd android && \
	mkdir -p app/libs && \
	gomobile bind -target=android -androidapi=21 -o app/libs/singbox.aar ../pkg/mobile
	@echo "构建 APK..."
	cd android && ./gradlew assembleRelease
	@echo "复制 APK 到 build 目录..."
	mkdir -p $(BUILD_DIR)
	cp android/app/build/outputs/apk/release/app-release-unsigned.apk $(BUILD_DIR)/singbox-xboard-android.apk
	@echo "Android APK 构建完成: $(BUILD_DIR)/singbox-xboard-android.apk"

# 构建所有平台
release: clean
	@echo "==> 构建所有平台版本..."
	@$(foreach platform,$(PLATFORMS),\
		echo "构建 $(platform)..." && \
		GOOS=$(word 1,$(subst /, ,$(platform))) \
		GOARCH=$(word 2,$(subst /, ,$(platform))) \
		$(GOBUILD) $(LDFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform)))$(if $(findstring windows,$(platform)),.exe,) \
		./cmd || exit 1;)
	@echo "所有平台构建完成"

# 运行测试
test:
	@echo "==> 运行测试..."
	$(GOTEST) -v ./...

# 清理
clean:
	@echo "==> 清理构建文件..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)

# 开发运行
run: build
	@echo "==> 运行 $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# 格式化代码
fmt:
	@echo "==> 格式化代码..."
	@gofmt -s -w .

# 代码检查
lint:
	@echo "==> 运行代码检查..."
	@golangci-lint run

# 下载 sing-box 核心
download-singbox:
	@echo "==> 下载 sing-box 核心..."
	@mkdir -p $(BUILD_DIR)/bin
	@./scripts/download-singbox.sh