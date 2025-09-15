# 多阶段构建 Dockerfile
FROM node:18-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制前端源代码
COPY frontend/ ./

# 构建前端
RUN npm run build

# Go 后端构建阶段
FROM golang:1.21-alpine AS backend-builder

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev sqlite-dev

# 设置工作目录
WORKDIR /app/backend

# 复制 go mod 文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码
COPY backend/ ./

# 构建二进制文件
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o flashcard cmd/server/main.go

# 最终运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 创建必要的目录
RUN mkdir -p /app/data /app/temp && \
    chown -R appuser:appgroup /app

# 从构建阶段复制文件
COPY --from=backend-builder /app/backend/flashcard /app/
COPY --from=frontend-builder /app/frontend/dist /app/static/

# 复制配置文件
COPY backend/env.example /app/.env

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release
ENV SQLITE_DB_PATH=/app/data/flashcard.db
ENV STATIC_PATH=/app/static

# 启动应用
CMD ["./flashcard"]
