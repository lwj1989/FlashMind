#!/bin/bash

# 检查是否在正确的目录中
if [ ! -f "test_data/go_flashcards.json" ]; then
    echo "请在FlashMind根目录运行此脚本"
    exit 1
fi

echo "启动FlashMind测试环境..."

# 启动后端服务器
echo "启动后端服务器..."
cd backend
go run main.go &
BACKEND_PID=$!
cd ..

# 等待后端服务器启动
echo "等待后端服务器启动..."
sleep 5

# 检查后端服务器是否成功启动
if ! curl -s http://localhost:8080/api/v1/decks > /dev/null; then
    echo "后端服务器启动失败"
    kill $BACKEND_PID
    exit 1
fi

# 启动前端服务器
echo "启动前端服务器..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

# 等待前端服务器启动
echo "等待前端服务器启动..."
sleep 10

# 检查前端服务器是否成功启动
if ! curl -s http://localhost:5173 > /dev/null; then
    echo "前端服务器启动失败"
    kill $BACKEND_PID $FRONTEND_PID
    exit 1
fi

# 导入测试数据
echo "导入测试数据..."
./import_test_data.sh

echo ""
echo "FlashMind测试环境已启动！"
echo "前端服务器: http://localhost:5173"
echo "后端服务器: http://localhost:8080"
echo ""
echo "按Ctrl+C停止服务器"

# 等待用户中断
trap "kill $BACKEND_PID $FRONTEND_PID; exit 0" INT
wait