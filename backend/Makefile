# Go Makefile for FlashMind Backend

.PHONY: help build run test test-coverage lint clean deps install

# 默认目标
help: ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# 构建
build: ## 构建应用
	@echo "构建后端应用..."
	@go build -o flashcard cmd/server/main.go

build-linux: ## 构建 Linux 版本
	@echo "构建 Linux 版本..."
	@CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o flashcard cmd/server/main.go

# 运行
run: ## 运行应用
	@echo "启动后端服务..."
	@go run cmd/server/main.go

# 测试
test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

test-race: ## 运行竞态检测测试
	@echo "运行竞态检测测试..."
	@go test -race -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试覆盖率..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

test-coverage-func: ## 显示函数级别的覆盖率
	@echo "函数级别覆盖率:"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

# 代码质量
lint: ## 运行代码检查
	@echo "运行代码检查..."
	@golangci-lint run

fmt: ## 格式化代码
	@echo "格式化代码..."
	@go fmt ./...

vet: ## 运行 go vet
	@echo "运行 go vet..."
	@go vet ./...

# 依赖管理
deps: ## 下载依赖
	@echo "下载依赖..."
	@go mod download

tidy: ## 整理依赖
	@echo "整理依赖..."
	@go mod tidy

vendor: ## 创建 vendor 目录
	@echo "创建 vendor 目录..."
	@go mod vendor

# 清理
clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -f flashcard
	@rm -f coverage.out coverage.html
	@rm -rf vendor/

# 安装工具
install-tools: ## 安装开发工具
	@echo "安装开发工具..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 数据库
db-reset: ## 重置数据库
	@echo "重置数据库..."
	@rm -f flashcard.db

# 开发
dev: deps ## 开发模式运行
	@echo "开发模式启动..."
	@air

# 生产环境
prod-build: test lint build-linux ## 生产环境构建

# 全面检查
check: fmt vet lint test-coverage ## 运行所有检查

# Docker
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	@docker build -t flashmind-backend .

docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	@docker run -p 8080:8080 flashmind-backend
