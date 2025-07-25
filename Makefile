# SingBox App Makefile

# 变量定义
APP_NAME = singbox-app
VERSION = 1.0.0
BUILD_DIR = build
MAIN_FILE = main.go

# Go编译参数
LDFLAGS = -ldflags "-X main.version=$(VERSION) -s -w"
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

# 默认目标
.PHONY: all
all: clean build

# 编译
.PHONY: build
build:
	@echo "正在编译 $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)"

# 交叉编译
.PHONY: build-linux
build-linux:
	@echo "正在为Linux编译..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_FILE)

.PHONY: build-windows
build-windows:
	@echo "正在为Windows编译..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_FILE)

.PHONY: build-darwin
build-darwin:
	@echo "正在为macOS编译..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_FILE)

.PHONY: build-all
build-all: build-linux build-windows build-darwin
	@echo "所有平台编译完成"

# 运行
.PHONY: run
run:
	go run $(MAIN_FILE)

.PHONY: run-web
run-web:
	go run $(MAIN_FILE) --web --port 8080

# 测试
.PHONY: test
test:
	go test ./...

# 清理
.PHONY: clean
clean:
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	go clean

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	go mod download
	go mod tidy

# 格式化代码
.PHONY: fmt
fmt:
	go fmt ./...

# 代码检查
.PHONY: vet
vet:
	go vet ./...

# 开发模式（热重载）
.PHONY: dev
dev:
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	air

# 帮助信息
.PHONY: help
help:
	@echo "可用命令:"
	@echo "  build       - 编译应用程序"
	@echo "  build-all   - 为所有平台编译"
	@echo "  run         - 运行应用程序（CLI模式）"
	@echo "  run-web     - 运行应用程序（Web模式）"
	@echo "  test        - 运行测试"
	@echo "  clean       - 清理构建文件"
	@echo "  deps        - 安装依赖"
	@echo "  fmt         - 格式化代码"
	@echo "  vet         - 代码检查"
	@echo "  dev         - 开发模式（热重载）"
	@echo "  help        - 显示此帮助信息"