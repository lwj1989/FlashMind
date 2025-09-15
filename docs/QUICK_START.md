# 快速开始指南

本指南将帮助您快速上手 FlashMind 智能记忆卡片系统。

## 🚀 5分钟快速启动

### 1. 获取代码
```bash
git clone https://github.com/lwj1989/FlashMind.git
cd FlashMind
```

### 2. 一键启动
```bash
./start.sh
```

### 3. 访问应用
- 前端界面: http://localhost:5173
- 后端API: http://localhost:8080

就这么简单！🎉

## 📝 第一次使用

### 创建您的第一个卡包

1. 点击 "卡包管理" 
2. 点击 "创建卡包"
3. 输入卡包名称，如 "英语单词"
4. 点击 "创建"

### 添加卡片

1. 进入刚创建的卡包
2. 点击 "添加卡片"
3. 输入问题和答案
4. 可选择添加标签分类
5. 支持 Markdown 格式和代码高亮

### 开始学习

1. 在卡包页面点击 "开始学习"
2. 选择学习模式
3. 根据记忆情况点击难度按钮
4. 系统会智能安排复习时间

## 🔄 快速导入数据

### 使用 TXT 格式
创建文件 `my_cards.txt`：
```
# 英语单词
apple
---
苹果
===
banana  
---
香蕉
```

然后在 "导入导出" 页面上传此文件。

### 使用测试数据
```bash
# 导入预置的 Go 语言学习卡片
cp test_data/go_flashcards.json ./
# 在导入页面上传 go_flashcards.json
```

## 🛠️ 常用命令

```bash
# 停止所有服务
./stop.sh

# 检查服务状态  
./status.sh

# 重启服务
./start.sh

# 查看日志
tail -f backend.log
tail -f frontend.log
```

## ❓ 遇到问题？

### 端口被占用
```bash
./stop.sh
./start.sh
```

### 服务启动失败
```bash
# 查看详细错误
./status.sh
cat backend.log
```

### 需要帮助
- 查看 [完整使用手册](../使用手册.md)
- 提交 [GitHub Issue](https://github.com/lwj1989/FlashMind/issues)

## 📚 下一步

- 了解 [高级功能](../使用手册.md#功能介绍)
- 查看 [API 文档](../使用手册.md#api文档)  
- 学习 [最佳实践](../使用手册.md#最佳实践)

开始您的高效学习之旅吧！ 🎓