#!/bin/bash

# 上传 singbox-xboard-client 到 GitHub 的脚本

echo "==================================="
echo "上传 singbox-xboard-client 到 GitHub"
echo "==================================="

# 检查是否在正确的目录
if [ ! -f "go.mod" ] || [ ! -d "cmd" ]; then
    echo "错误：请在 singbox-xboard-client 目录下运行此脚本"
    exit 1
fi

# 初始化 Git 仓库
if [ ! -d ".git" ]; then
    echo "初始化 Git 仓库..."
    git init
fi

# 添加所有文件
echo "添加所有文件到 Git..."
git add .

# 创建初始提交
echo "创建初始提交..."
git commit -m "Initial commit: singbox-xboard-client with full platform support

Features:
- sing-box core integration
- xboard panel support
- Web UI management interface
- Multi-platform support (Windows, macOS, Linux, Android)
- GitHub Actions automated builds
- Support for latest protocols: Hysteria2, VLESS+Reality, etc."

# 提示用户输入仓库信息
echo ""
echo "请按照以下步骤操作："
echo ""
echo "1. 在 GitHub 上创建一个新的仓库（如果还没有创建）"
echo "   - 仓库名建议：singbox-xboard-client"
echo "   - 设置为 Public（如果你想让其他人也能使用）"
echo "   - 不要初始化 README、.gitignore 或 LICENSE"
echo ""
echo "2. 复制下面的命令并替换 YOUR_USERNAME 为你的 GitHub 用户名："
echo ""
echo "   git remote add origin https://github.com/YOUR_USERNAME/singbox-xboard-client.git"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "3. 如果你想触发自动构建，运行："
echo "   git tag v1.0.0"
echo "   git push origin v1.0.0"
echo ""
echo "4. 稍等几分钟后，在 GitHub Actions 页面查看构建进度"
echo "   构建完成后，在 Releases 页面下载各平台的客户端"
echo ""

# 可选：自动设置远程仓库（需要用户输入）
read -p "是否现在设置 GitHub 远程仓库？(y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "请输入你的 GitHub 用户名: " username
    if [ ! -z "$username" ]; then
        git remote add origin "https://github.com/$username/singbox-xboard-client.git"
        echo "远程仓库已设置为: https://github.com/$username/singbox-xboard-client.git"
        
        read -p "是否现在推送到 GitHub？(y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git branch -M main
            echo "正在推送到 GitHub..."
            git push -u origin main
            
            read -p "是否创建 v1.0.0 标签并触发自动构建？(y/n): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                git tag v1.0.0
                git push origin v1.0.0
                echo "已创建标签 v1.0.0，GitHub Actions 将开始自动构建"
            fi
        fi
    fi
fi

echo ""
echo "完成！"