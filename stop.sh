#!/bin/bash

# FlashMind 一键停止脚本
# 作者: FlashMind Team
# 描述: 停止前后端服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# PID 文件路径
BACKEND_PID_FILE="$PROJECT_ROOT/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_ROOT/.frontend.pid"

echo -e "${BLUE}===========================================${NC}"
echo -e "${BLUE}        FlashMind 停止脚本 v1.0           ${NC}"
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

    local stopped=false

    # 通过 PID 文件停止
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null 2>&1; then
            print_info "发现 $service_name 进程 (PID: $pid)，正在终止..."
            kill $pid 2>/dev/null || true
            
            # 等待进程优雅关闭
            local count=0
            while ps -p $pid > /dev/null 2>&1 && [ $count -lt 10 ]; do
                sleep 1
                count=$((count + 1))
            done
            
            # 如果进程仍然存在，强制终止
            if ps -p $pid > /dev/null 2>&1; then
                print_warning "进程未响应，强制终止 $service_name 进程..."
                kill -9 $pid 2>/dev/null || true
                sleep 1
            fi
            stopped=true
        fi
        rm -f "$pid_file"
    fi

    # 通过端口停止（备用方案）
    if check_port $port; then
        print_info "发现端口 $port 仍被占用，正在终止相关进程..."
        local pids=$(lsof -ti:$port)
        if [ -n "$pids" ]; then
            echo "$pids" | xargs kill 2>/dev/null || true
            sleep 2
            
            # 检查是否还有进程占用端口
            if check_port $port; then
                print_warning "强制终止端口 $port 上的进程..."
                lsof -ti:$port | xargs kill -9 2>/dev/null || true
                sleep 1
            fi
        fi
        stopped=true
    fi

    if $stopped; then
        print_success "$service_name 服务已停止"
    else
        print_info "$service_name 服务未运行"
    fi
}

# 函数：清理临时文件
cleanup() {
    print_info "正在清理临时文件..."
    
    # 清理 PID 文件
    rm -f "$BACKEND_PID_FILE" "$FRONTEND_PID_FILE"
    
    # 清理日志文件（可选）
    if [ "$1" = "--clean-logs" ]; then
        print_info "正在清理日志文件..."
        rm -f "$PROJECT_ROOT/backend.log" "$PROJECT_ROOT/frontend.log"
        print_success "日志文件已清理"
    fi
}

# 函数：检查服务状态
check_status() {
    print_info "正在检查服务状态..."
    
    local any_running=false
    
    # 检查后端
    if check_port 8080; then
        print_warning "端口 8080 仍被占用"
        any_running=true
    fi
    
    # 检查前端
    if check_port 5173; then
        print_warning "端口 5173 仍被占用"
        any_running=true
    fi
    
    if ! $any_running; then
        print_success "所有服务已完全停止"
    else
        print_warning "部分端口仍被占用，可能有其他程序在使用"
    fi
}

# 显示帮助信息
show_help() {
    echo "FlashMind 停止脚本"
    echo ""
    echo "用法:"
    echo "  $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --clean-logs    同时清理日志文件"
    echo "  --help          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                # 停止所有服务"
    echo "  $0 --clean-logs  # 停止所有服务并清理日志"
}

# 主逻辑
main() {
    # 解析命令行参数
    local clean_logs=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --clean-logs)
                clean_logs=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "开始停止 FlashMind 应用..."
    
    # 停止服务
    stop_service "前端" "$FRONTEND_PID_FILE" 5173
    stop_service "后端" "$BACKEND_PID_FILE" 8080
    
    # 清理文件
    if $clean_logs; then
        cleanup --clean-logs
    else
        cleanup
    fi
    
    echo ""
    
    # 检查状态
    check_status
    
    echo ""
    print_success "FlashMind 应用已停止！"
    echo ""
    echo -e "${YELLOW}提示:${NC}"
    echo -e "  - 使用 ${GREEN}./start.sh${NC} 启动所有服务"
    echo -e "  - 使用 ${GREEN}./stop.sh --clean-logs${NC} 停止并清理日志文件"
    echo ""
}

# 执行主逻辑
main "$@"
