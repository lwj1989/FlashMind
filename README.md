# FlashMind - 智能记忆卡片系统

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/your-username/FlashMind)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.0+-4FC08D.svg)](https://vuejs.org/)

FlashMind 是一个基于间隔重复记忆算法的智能学习系统，帮助用户高效地学习和记忆各种知识点。

## ✨ 核心特性

- 🧠 **智能复习算法** - 基于艾宾浩斯遗忘曲线的间隔重复算法
- 📚 **多格式支持** - 支持文本、Markdown、代码片段等多种格式
- 🏷️ **灵活标签系统** - 多层级标签管理，便于分类和检索
- 📊 **学习统计** - 详细的学习进度追踪和数据分析
- 🔄 **批量导入导出** - 支持 TXT、CSV、JSON 等多种格式
- 🎯 **个性化学习** - 多种学习模式（复习、随机、指定卡包）
- 📱 **响应式设计** - 支持桌面端和移动端访问

## 🚀 快速开始

### 系统要求

- **Go**: 1.21 或更高版本
- **Node.js**: 16.0 或更高版本
- **npm**: 8.0 或更高版本

### 一键启动

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd FlashMind
   ```

2. **启动所有服务**
   ```bash
   ./start.sh
   ```
   
   这个命令会：
   - 自动安装所有依赖
   - 编译后端服务
   - 启动前后端服务
   - 如果服务已运行，会先停止再启动（重启功能）

3. **访问应用**
   - 前端界面: http://localhost:5173
   - 后端API: http://localhost:8080

### 管理服务

```bash
# 停止所有服务
./stop.sh

# 停止服务并清理日志
./stop.sh --clean-logs

# 检查服务状态
./status.sh

# 重启服务（等同于 start.sh）
./start.sh
```

## 📖 使用指南

### 基本操作

1. **创建卡包**: 在卡包管理页面创建新的学习卡包
2. **添加标签**: 为卡片分类创建标签系统
3. **创建卡片**: 添加问题和答案，支持 Markdown 格式
4. **开始学习**: 选择学习模式进行复习

### 快速导入数据

支持 TXT 格式的快速导入：

```txt
# 卡包：Go语言基础
# 标签：语法,基础

什么是Go语言？
Go是Google开发的开源编程语言，具有简洁、高效、并发性强的特点。

如何声明变量？
使用var关键字或:=操作符。例如：var name string 或 name := "hello"

---

什么是goroutine？
goroutine是Go语言的轻量级线程，通过go关键字启动。
```

## 📁 项目结构

```
FlashMind/
├── backend/                # Go后端
│   ├── cmd/server/        # 主程序入口
│   ├── internal/          # 内部包
│   │   ├── handlers/      # HTTP处理器
│   │   ├── models/        # 数据模型
│   │   ├── services/      # 业务逻辑
│   │   └── middleware/    # 中间件
│   └── pkg/               # 公共包
├── frontend/              # Vue前端
│   ├── src/
│   │   ├── components/    # Vue组件
│   │   ├── views/         # 页面视图
│   │   ├── api/           # API调用
│   │   └── router/        # 路由配置
│   └── public/            # 静态资源
├── start.sh               # 一键启动脚本
├── stop.sh                # 一键停止脚本
├── status.sh              # 状态检查脚本
└── README.md              # 项目说明
```

## 🔧 开发指南

### 手动启动（开发模式）

**后端开发**:
```bash
cd backend
go mod download
go run cmd/server/main.go
```

**前端开发**:
```bash
cd frontend
npm install
npm run dev
```

### 构建生产版本

**后端**:
```bash
cd backend
go build -o flashcard cmd/server/main.go
```

**前端**:
```bash
cd frontend
npm run build
```

### 运行测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm run test
```

## 📊 API 接口

### 核心接口

- `GET /api/decks` - 获取卡包列表
- `POST /api/decks` - 创建卡包
- `GET /api/cards` - 获取卡片列表
- `POST /api/cards` - 创建卡片
- `GET /api/tags` - 获取标签列表
- `POST /api/study/start` - 开始学习会话
- `POST /api/import` - 批量导入数据

详细的 API 文档请参考 [使用手册.md](使用手册.md)

## 📝 日志和调试

- **后端日志**: `backend.log`
- **前端日志**: `frontend.log`

实时查看日志：
```bash
# 后端日志
tail -f backend.log

# 前端日志
tail -f frontend.log
```

## 🔍 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 检查端口占用
   lsof -i :8080
   lsof -i :5173
   
   # 使用脚本停止服务
   ./stop.sh
   ```

2. **服务启动失败**
   ```bash
   # 检查服务状态
   ./status.sh
   
   # 查看详细日志
   cat backend.log
   cat frontend.log
   ```

3. **依赖安装失败**
   ```bash
   # 清理并重新安装
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   
   cd ../backend
   go clean -modcache
   go mod download
   ```

### 重置环境

完全重置开发环境：
```bash
./stop.sh --clean-logs
rm -f .backend.pid .frontend.pid
rm -f backend.log frontend.log
rm -f backend/flashcard.db
./start.sh
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🔗 相关链接

- [详细使用手册](使用手册.md)
- [API文档](使用手册.md#api文档)
- [TXT导入格式说明](使用手册.md#txt格式说明)
- [常见问题解答](使用手册.md#常见问题解答)

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 邮箱: flashmind@example.com
- GitHub Issues: [创建 Issue](https://github.com/your-username/FlashMind/issues)

---

**Happy Learning! 🎓**
