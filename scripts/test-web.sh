#!/bin/bash

# SingBox App Web UI 测试脚本

echo "=== SingBox App Web UI 测试 ==="

# 检查编译结果
if [ ! -f "./build/singbox-app" ]; then
    echo "❌ 应用程序未编译，请先运行: make build"
    exit 1
fi

echo "✅ 应用程序已编译"

# 启动Web界面
echo ""
echo "🌐 启动Web界面模式:"
echo "命令: ./build/singbox-app --web --port 8080"
echo ""
echo "📱 访问地址: http://localhost:8080"
echo "按Ctrl+C停止服务"
echo ""

# 启动应用程序
./build/singbox-app --web --port 8080