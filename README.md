# FlashMind 📚

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.0+-green.svg)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

FlashMind 是一个基于间隔重复记忆算法的智能学习系统，帮助用户高效地学习和记忆各种知识点。

## ✨ 核心特性

- 🧠 **智能复习算法** - 基于艾宾浩斯遗忘曲线的间隔重复算法
- 📚 **多格式支持** - 支持 Markdown、代码高亮、数学公式等多种格式
- 🏷️ **灵活标签系统** - 多层级标签管理，便于分类和检索
- 📊 **学习统计** - 详细的学习进度追踪和数据分析
- 🔄 **批量导入导出** - 支持 TXT、CSV、JSON 等多种格式
- 🎯 **个性化学习** - 多种学习模式（复习、随机、指定卡包）
- 📱 **响应式设计** - 支持桌面端和移动端访问

## 🚀 快速开始

### 系统要求

- **Go**: 1.21+ 
- **Node.js**: 16.0+
- **npm**: 8.0+

### 一键启动

```bash
# 克隆项目
git clone https://github.com/lwj1989/FlashMind.git
cd FlashMind

# 启动所有服务
./start.sh
```

启动成功后访问：
- **前端界面**: http://localhost:5173
- **后端API**: http://localhost:8080

### 服务管理

```bash
# 停止所有服务
./stop.sh

# 检查服务状态
./status.sh

# 重启服务
./start.sh
```

## 📖 使用指南

### 基本操作

1. **创建卡包**: 在卡包管理页面创建新的学习卡包
2. **添加标签**: 为卡片分类创建标签系统
3. **创建卡片**: 添加问题和答案，支持 Markdown 格式
4. **开始学习**: 选择学习模式进行复习

### 数据导入

支持多种格式的数据导入：

#### TXT 格式
```
# 标签名称
问题1
---
答案1
===
问题2
---
答案2
```

#### CSV 格式
```csv
Question,Answer,Tag
问题1,答案1,标签1
问题2,答案2,标签2
```

#### JSON 格式
```json
{
  "name": "卡包名称",
  "cards": [
    {
      "question": "问题1",
      "answer": "答案1",
      "tag_name": "标签1"
    }
  ]
}
```

## 🏗️ 项目架构

### 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue 3 + Vite + Tailwind CSS + Element Plus
- **数据库**: SQLite（可扩展到 PostgreSQL/MySQL）

### 项目结构

```
FlashMind/
├── backend/                    # Go 后端服务
│   ├── cmd/server/            # 主程序入口
│   ├── internal/              # 内部包
│   │   ├── handlers/          # HTTP 处理器
│   │   ├── models/            # 数据模型
│   │   ├── services/          # 业务逻辑
│   │   ├── middleware/        # 中间件
│   │   └── config/            # 配置管理
│   └── pkg/                   # 公共包
├── frontend/                  # Vue 前端应用
│   ├── src/
│   │   ├── components/        # Vue 组件
│   │   ├── views/             # 页面视图
│   │   ├── api/               # API 调用
│   │   └── router/            # 路由配置
│   └── public/                # 静态资源
├── test_data/                 # 测试数据
├── start.sh                   # 启动脚本
├── stop.sh                    # 停止脚本
└── status.sh                  # 状态检查脚本
```

## 🛠️ 开发指南

### 本地开发

#### 后端开发
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### 前端开发
```bash
cd frontend
npm install
npm run dev
```

### 构建部署

#### 生产构建
```bash
# 后端构建
cd backend
go build -o flashcard cmd/server/main.go

# 前端构建
cd frontend
npm run build
```

### 测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm run test
```

## 📋 API 接口

### 核心接口

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/v1/decks` | GET | 获取卡包列表 |
| `/api/v1/decks` | POST | 创建卡包 |
| `/api/v1/cards` | GET | 获取卡片列表 |
| `/api/v1/cards` | POST | 创建卡片 |
| `/api/v1/tags` | GET | 获取标签列表 |
| `/api/v1/study/start` | POST | 开始学习会话 |
| `/api/v1/import-export/decks` | POST | 导入数据 |

完整的 API 文档请参考 [使用手册](使用手册.md)

## 🔍 故障排除

### 常见问题

**端口被占用**
```bash
# 检查端口占用
lsof -i :8080
lsof -i :5173

# 停止服务
./stop.sh
```

**服务启动失败**
```bash
# 检查服务状态
./status.sh

# 查看日志
tail -f backend.log
tail -f frontend.log
```

**依赖安装失败**
```bash
# 清理并重新安装
cd frontend
rm -rf node_modules package-lock.json
npm install

cd ../backend
go clean -modcache
go mod download
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

### 开发规范

- 代码提交前请运行测试确保通过
- 遵循现有的代码风格和命名规范
- 添加适当的注释和文档
- 确保新功能有相应的测试覆盖

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📞 联系方式

- **项目主页**: https://github.com/lwj1989/FlashMind
- **问题反馈**: [创建 Issue](https://github.com/lwj1989/FlashMind/issues)
- **功能建议**: [讨论区](https://github.com/lwj1989/FlashMind/discussions)

---

**Happy Learning! 🎓**