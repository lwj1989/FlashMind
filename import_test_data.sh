#!/bin/bash

# 检查后端服务器是否运行
if ! curl -s http://localhost:8080/api/v1/decks > /dev/null; then
    echo "后端服务器未运行，请先启动后端服务器"
    exit 1
fi

# 检查前端服务器是否运行
if ! curl -s http://localhost:5173 > /dev/null; then
    echo "前端服务器未运行，请先启动前端服务器"
    exit 1
fi

echo "开始导入测试数据..."

# 使用curl导入测试数据
curl -X POST \
  -H "Content-Type: multipart/form-data" \
  -F "file=@test_data/go_flashcards.json" \
  http://localhost:8080/api/v1/import/deck

echo ""
echo "测试数据导入完成！"
echo "请访问 http://localhost:5173 查看导入的卡包和卡片"