#!/bin/bash

# 进入前端目录
cd frontend

# 安装依赖
echo "安装前端测试依赖..."
npm install

# 运行测试
echo "运行前端测试..."
npm run test:coverage

# 返回根目录
cd ..

# 进入后端目录
cd backend

# 安装依赖
echo "安装后端测试依赖..."
go mod tidy

# 运行测试
echo "运行后端测试..."
go test -v -cover ./...

# 返回根目录
cd ..

echo "所有测试完成！"