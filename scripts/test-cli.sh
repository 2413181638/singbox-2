#!/bin/bash

# SingBox App CLI 测试脚本

echo "=== SingBox App CLI 测试 ==="

# 检查编译结果
if [ ! -f "./build/singbox-app" ]; then
    echo "❌ 应用程序未编译，请先运行: make build"
    exit 1
fi

echo "✅ 应用程序已编译"

# 显示帮助信息
echo ""
echo "📖 显示帮助信息:"
./build/singbox-app --help

echo ""
echo "🚀 启动CLI模式 (按Ctrl+C停止):"
echo "命令: ./build/singbox-app"
echo ""

# 启动应用程序
./build/singbox-app