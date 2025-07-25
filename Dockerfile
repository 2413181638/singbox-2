# 多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git make

# 复制源代码
COPY . .

# 下载Go模块
RUN go mod download

# 编译应用
RUN make build

# 运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -S singbox && adduser -S singbox -G singbox

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/build/singbox-app .
COPY --from=builder /app/web ./web
COPY --from=builder /app/examples ./examples

# 创建配置目录
RUN mkdir -p /etc/singbox-app && \
    chown -R singbox:singbox /app /etc/singbox-app

# 复制默认配置
COPY examples/config-socks.yaml /etc/singbox-app/config.yaml

# 切换到非root用户
USER singbox

# 暴露端口
EXPOSE 1080 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/status || exit 1

# 默认启动Web界面
CMD ["./singbox-app", "--web", "--port", "8080", "--config", "/etc/singbox-app/config.yaml"]