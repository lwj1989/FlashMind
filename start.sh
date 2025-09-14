#!/bin/bash

# FlashMind 一键启动脚本
# 作者: FlashMind Team
# 描述: 启动前后端服务，支持重启功能

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# PID 文件路径
BACKEND_PID_FILE="$PROJECT_ROOT/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_ROOT/.frontend.pid"

echo -e "${BLUE}===========================================${NC}"
echo -e "${BLUE}        FlashMind 启动脚本 v1.0           ${NC}"
echo -e "${BLUE}===========================================${NC}"

# 函数：打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 函数：检查端口是否被占用
check_port() {
    local port=$1
    if lsof -ti:$port > /dev/null 2>&1; then
        return 0  # 端口被占用
    else
        return 1  # 端口未被占用
    fi
}

# 函数：停止服务
stop_service() {
    local service_name=$1
    local pid_file=$2
    local port=$3

    print_info "正在停止 $service_name 服务..."

    # 通过 PID 文件停止
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null 2>&1; then
            print_info "发现 $service_name 进程 (PID: $pid)，正在终止..."
            kill $pid 2>/dev/null || true
            sleep 2
            
            # 如果进程仍然存在，强制终止
            if ps -p $pid > /dev/null 2>&1; then
                print_warning "强制终止 $service_name 进程..."
                kill -9 $pid 2>/dev/null || true
            fi
        fi
        rm -f "$pid_file"
    fi

    # 通过端口停止（备用方案）
    if check_port $port; then
        print_info "发现端口 $port 仍被占用，正在终止相关进程..."
        lsof -ti:$port | xargs kill -9 2>/dev/null || true
        sleep 1
    fi

    print_success "$service_name 服务已停止"
}

# 函数：启动后端服务
start_backend() {
    print_info "正在启动后端服务..."
    
    cd "$BACKEND_DIR"
    
    # 检查 Go 环境
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go 语言环境"
        exit 1
    fi

    # 检查依赖
    if [ ! -f "go.mod" ]; then
        print_error "未找到 go.mod 文件"
        exit 1
    fi

    # 下载依赖
    print_info "正在下载 Go 依赖..."
    go mod download

    # 编译并启动
    print_info "正在编译后端服务..."
    go build -o flashcard cmd/server/main.go

    if [ ! -f "./flashcard" ]; then
        print_error "后端编译失败"
        exit 1
    fi

    print_info "正在启动后端服务 (端口: 8080)..."
    nohup ./flashcard > ../backend.log 2>&1 &
    echo $! > "$BACKEND_PID_FILE"
    
    # 等待服务启动
    sleep 3
    
    # 检查服务是否成功启动
    if ps -p $(cat "$BACKEND_PID_FILE") > /dev/null 2>&1; then
        print_success "后端服务启动成功 (PID: $(cat "$BACKEND_PID_FILE"))"
        print_info "后端服务地址: http://localhost:8080"
        print_info "日志文件: $PROJECT_ROOT/backend.log"
    else
        print_error "后端服务启动失败，请检查日志文件"
        cat "$PROJECT_ROOT/backend.log" 2>/dev/null || true
        exit 1
    fi
}

# 函数：启动前端服务
start_frontend() {
    print_info "正在启动前端服务..."
    
    cd "$FRONTEND_DIR"
    
    # 检查 Node.js 环境
    if ! command -v npm &> /dev/null; then
        print_error "npm 未安装，请先安装 Node.js 环境"
        exit 1
    fi

    # 检查依赖
    if [ ! -f "package.json" ]; then
        print_error "未找到 package.json 文件"
        exit 1
    fi

    # 安装依赖
    if [ ! -d "node_modules" ]; then
        print_info "正在安装前端依赖..."
        npm install
    fi

    print_info "正在启动前端开发服务器 (端口: 5173)..."
    nohup npm run dev > ../frontend.log 2>&1 &
    echo $! > "$FRONTEND_PID_FILE"
    
    # 等待服务启动
    sleep 5
    
    # 检查服务是否成功启动
    if ps -p $(cat "$FRONTEND_PID_FILE") > /dev/null 2>&1; then
        print_success "前端服务启动成功 (PID: $(cat "$FRONTEND_PID_FILE"))"
        print_info "前端服务地址: http://localhost:5173"
        print_info "日志文件: $PROJECT_ROOT/frontend.log"
    else
        print_error "前端服务启动失败，请检查日志文件"
        cat "$PROJECT_ROOT/frontend.log" 2>/dev/null || true
        exit 1
    fi
}

# 函数：检查服务状态
check_status() {
    print_info "正在检查服务状态..."
    
    local backend_running=false
    local frontend_running=false
    
    # 检查后端
    if [ -f "$BACKEND_PID_FILE" ] && ps -p $(cat "$BACKEND_PID_FILE") > /dev/null 2>&1; then
        backend_running=true
        print_success "后端服务运行中 (PID: $(cat "$BACKEND_PID_FILE"))"
    elif check_port 8080; then
        print_warning "后端端口 8080 被占用，但非本程序启动"
    else
        print_info "后端服务未运行"
    fi
    
    # 检查前端
    if [ -f "$FRONTEND_PID_FILE" ] && ps -p $(cat "$FRONTEND_PID_FILE") > /dev/null 2>&1; then
        frontend_running=true
        print_success "前端服务运行中 (PID: $(cat "$FRONTEND_PID_FILE"))"
    elif check_port 5173; then
        print_warning "前端端口 5173 被占用，但非本程序启动"
    else
        print_info "前端服务未运行"
    fi
    
    if $backend_running && $frontend_running; then
        print_success "所有服务运行正常"
        echo ""
        echo -e "${GREEN}访问地址:${NC}"
        echo -e "  前端: ${BLUE}http://localhost:5173${NC}"
        echo -e "  后端: ${BLUE}http://localhost:8080${NC}"
    fi
}

# 主逻辑
main() {
    print_info "开始启动 FlashMind 应用..."
    
    # 停止现有服务（重启逻辑）
    stop_service "后端" "$BACKEND_PID_FILE" 8080
    stop_service "前端" "$FRONTEND_PID_FILE" 5173
    
    # 等待端口释放
    sleep 2
    
    # 启动服务
    start_backend
    start_frontend
    
    echo ""
    print_success "FlashMind 应用启动完成！"
    echo ""
    
    # 显示状态
    check_status
    
    echo ""
    echo -e "${YELLOW}提示:${NC}"
    echo -e "  - 使用 ${GREEN}./stop.sh${NC} 停止所有服务"
    echo -e "  - 使用 ${GREEN}./start.sh${NC} 重启所有服务"
    echo -e "  - 查看日志: ${GREEN}tail -f backend.log${NC} 或 ${GREEN}tail -f frontend.log${NC}"
    echo ""
}

# 执行主逻辑
main
