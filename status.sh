#!/bin/bash

# FlashMind 状态检查脚本
# 作者: FlashMind Team
# 描述: 检查前后端服务运行状态

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
echo -e "${BLUE}        FlashMind 状态检查 v1.0           ${NC}"
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

# 函数：检查HTTP服务
check_http_service() {
    local url=$1
    local name=$2
    
    if curl -s --connect-timeout 3 "$url" > /dev/null 2>&1; then
        print_success "$name HTTP服务响应正常"
        return 0
    else
        print_error "$name HTTP服务无响应"
        return 1
    fi
}

# 函数：获取进程信息
get_process_info() {
    local pid=$1
    if ps -p $pid -o pid,ppid,pcpu,pmem,etime,command --no-headers 2>/dev/null; then
        return 0
    else
        return 1
    fi
}

# 函数：检查服务详细状态
check_service_detail() {
    local service_name=$1
    local pid_file=$2
    local port=$3
    local http_url=$4
    
    echo ""
    echo -e "${BLUE}=== $service_name 服务状态 ===${NC}"
    
    local service_running=false
    local port_occupied=false
    local http_ok=false
    
    # 检查PID文件
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        echo -e "PID文件: ${GREEN}存在${NC} ($pid_file)"
        echo -e "PID: $pid"
        
        if ps -p $pid > /dev/null 2>&1; then
            print_success "进程运行中"
            service_running=true
            echo "进程详情:"
            get_process_info $pid | while read line; do
                echo "  $line"
            done
        else
            print_error "PID文件存在但进程不存在"
        fi
    else
        echo -e "PID文件: ${RED}不存在${NC}"
    fi
    
    # 检查端口
    if check_port $port; then
        print_success "端口 $port 被占用"
        port_occupied=true
        local port_pid=$(lsof -ti:$port)
        echo "占用端口的进程: $port_pid"
    else
        print_error "端口 $port 未被占用"
    fi
    
    # 检查HTTP服务
    if [ -n "$http_url" ]; then
        if check_http_service "$http_url" "$service_name"; then
            http_ok=true
        fi
    fi
    
    # 综合状态
    if $service_running && $port_occupied && $http_ok; then
        echo -e "综合状态: ${GREEN}正常运行${NC}"
    elif $port_occupied; then
        echo -e "综合状态: ${YELLOW}端口被占用但可能非本程序${NC}"
    else
        echo -e "综合状态: ${RED}未运行${NC}"
    fi
}

# 函数：显示系统信息
show_system_info() {
    echo ""
    echo -e "${BLUE}=== 系统信息 ===${NC}"
    echo "操作系统: $(uname -s)"
    echo "内核版本: $(uname -r)"
    echo "主机名: $(hostname)"
    echo "当前时间: $(date)"
    echo "运行时长: $(uptime | awk '{print $3, $4}' | sed 's/,//')"
    
    # 内存使用情况
    if command -v free &> /dev/null; then
        echo "内存使用:"
        free -h | grep -E "(Mem|Swap):"
    fi
    
    # 磁盘使用情况
    echo "磁盘使用:"
    df -h "$PROJECT_ROOT" | tail -1
}

# 函数：显示日志摘要
show_log_summary() {
    echo ""
    echo -e "${BLUE}=== 日志摘要 ===${NC}"
    
    local backend_log="$PROJECT_ROOT/backend.log"
    local frontend_log="$PROJECT_ROOT/frontend.log"
    
    # 后端日志
    if [ -f "$backend_log" ]; then
        echo "后端日志 ($backend_log):"
        echo "  文件大小: $(du -h "$backend_log" | cut -f1)"
        echo "  最后修改: $(stat -c %y "$backend_log" 2>/dev/null || stat -f %m "$backend_log" 2>/dev/null || echo "未知")"
        if [ -s "$backend_log" ]; then
            echo "  最后几行:"
            tail -3 "$backend_log" | sed 's/^/    /'
        else
            echo "  日志文件为空"
        fi
    else
        echo "后端日志: 不存在"
    fi
    
    echo ""
    
    # 前端日志
    if [ -f "$frontend_log" ]; then
        echo "前端日志 ($frontend_log):"
        echo "  文件大小: $(du -h "$frontend_log" | cut -f1)"
        echo "  最后修改: $(stat -c %y "$frontend_log" 2>/dev/null || stat -f %m "$frontend_log" 2>/dev/null || echo "未知")"
        if [ -s "$frontend_log" ]; then
            echo "  最后几行:"
            tail -3 "$frontend_log" | sed 's/^/    /'
        else
            echo "  日志文件为空"
        fi
    else
        echo "前端日志: 不存在"
    fi
}

# 函数：显示网络信息
show_network_info() {
    echo ""
    echo -e "${BLUE}=== 网络信息 ===${NC}"
    
    # 检查网络连接
    echo "网络连接状态:"
    if netstat -tuln 2>/dev/null | grep -E ":(8080|5173)" | head -10; then
        true
    else
        echo "  无相关端口监听"
    fi
    
    # 本机IP地址
    echo ""
    echo "本机网络接口:"
    if command -v ip &> /dev/null; then
        ip addr show | grep -E "inet " | grep -v "127.0.0.1" | head -3
    elif command -v ifconfig &> /dev/null; then
        ifconfig | grep -E "inet " | grep -v "127.0.0.1" | head -3
    else
        echo "  无法获取网络接口信息"
    fi
}

# 主逻辑
main() {
    print_info "正在检查 FlashMind 应用状态..."
    
    # 检查后端服务
    check_service_detail "后端" "$BACKEND_PID_FILE" 8080 "http://localhost:8080/api/health"
    
    # 检查前端服务
    check_service_detail "前端" "$FRONTEND_PID_FILE" 5173 "http://localhost:5173"
    
    # 显示系统信息
    show_system_info
    
    # 显示日志摘要
    show_log_summary
    
    # 显示网络信息
    show_network_info
    
    echo ""
    echo -e "${BLUE}===========================================${NC}"
    print_success "状态检查完成"
    echo ""
    echo -e "${YELLOW}提示:${NC}"
    echo -e "  - 使用 ${GREEN}./start.sh${NC} 启动所有服务"
    echo -e "  - 使用 ${GREEN}./stop.sh${NC} 停止所有服务"
    echo -e "  - 查看实时日志: ${GREEN}tail -f backend.log${NC} 或 ${GREEN}tail -f frontend.log${NC}"
    echo ""
}

# 执行主逻辑
main "$@"
