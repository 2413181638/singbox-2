#!/bin/bash

# SingBox App 多平台编译脚本

echo "=== SingBox App 多平台编译 ==="

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请先安装Go"
    exit 1
fi

echo "✅ Go环境检查通过"

# 清理之前的构建
echo "🧹 清理之前的构建..."
make clean

# 安装依赖
echo "📦 安装依赖..."
make deps

# 编译所有平台
echo "🔨 开始多平台编译..."
make build-all

# 检查编译结果
echo ""
echo "📋 编译结果:"
ls -la build/

echo ""
echo "✅ 多平台编译完成!"
echo ""
echo "📁 可执行文件位置:"
echo "  - Linux:   build/singbox-app-linux-amd64"
echo "  - Windows: build/singbox-app-windows-amd64.exe"
echo "  - macOS:   build/singbox-app-darwin-amd64"