#!/bin/bash

# SingBox App Docker 运行脚本

echo "=== SingBox App Docker 部署 ==="

# 检查Docker环境
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

echo "✅ Docker环境检查通过"

# 创建必要的目录
echo "📁 创建配置和日志目录..."
mkdir -p config logs

# 复制示例配置
if [ ! -f "config/config.yaml" ]; then
    echo "📋 复制示例配置文件..."
    cp examples/config-socks.yaml config/config.yaml
fi

# 构建并启动服务
echo "🐳 构建并启动Docker服务..."
docker-compose up --build -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

echo ""
echo "✅ SingBox App已启动!"
echo ""
echo "🌐 Web管理界面: http://localhost:8080"
echo "🔌 SOCKS代理端口: localhost:1080"
echo ""
echo "📋 常用命令:"
echo "  查看日志: docker-compose logs -f"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart"
echo ""