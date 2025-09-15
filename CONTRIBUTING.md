# 贡献指南

感谢您对 FlashMind 项目的关注和贡献！

## 🚀 如何贡献

### 报告问题
- 使用 [GitHub Issues](https://github.com/lwj1989/FlashMind/issues) 报告 Bug
- 提供详细的问题描述和复现步骤
- 包含您的环境信息（操作系统、Go版本、Node.js版本等）

### 提交功能请求
- 在 Issues 中详细描述您的需求
- 说明功能的使用场景和价值
- 欢迎提供设计方案或实现思路

### 提交代码

#### 1. 准备工作
```bash
# Fork 项目到您的 GitHub 账户
# 克隆您的 Fork
git clone https://github.com/YOUR_USERNAME/FlashMind.git
cd FlashMind

# 添加上游仓库
git remote add upstream https://github.com/lwj1989/FlashMind.git
```

#### 2. 创建功能分支
```bash
# 更新主分支
git checkout main
git pull upstream main

# 创建新分支
git checkout -b feature/your-feature-name
```

#### 3. 开发和测试
```bash
# 安装依赖
./start.sh

# 运行测试
cd backend && go test ./...
cd ../frontend && npm test
```

#### 4. 提交更改
```bash
git add .
git commit -m "feat: add amazing feature"
git push origin feature/your-feature-name
```

#### 5. 创建 Pull Request
- 在 GitHub 上创建 Pull Request
- 提供清晰的标题和描述
- 关联相关的 Issues

## 📝 开发规范

### 代码风格

#### Go 代码
- 使用 `go fmt` 格式化代码
- 遵循 Go 官方代码规范
- 添加适当的注释，特别是公共函数和结构体
- 使用有意义的变量和函数名

#### JavaScript/Vue 代码
- 使用 ESLint 和 Prettier 保持代码风格一致
- 遵循 Vue 3 最佳实践
- 组件名使用 PascalCase
- 文件名使用 kebab-case

### 提交信息规范

使用 [Conventional Commits](https://conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Type 类型：**
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或工具的变动

**示例：**
```
feat(backend): add card search functionality
fix(frontend): resolve card display issue on mobile
docs: update installation guide
```

### 测试要求

- 新功能必须包含相应的测试
- 确保所有现有测试通过
- 测试覆盖率应保持在合理水平

#### 后端测试
```bash
cd backend
go test ./... -v
go test ./... -cover
```

#### 前端测试
```bash
cd frontend
npm test
npm run test:coverage
```

## 🏗️ 项目结构

### 后端结构
```
backend/
├── cmd/server/          # 主程序入口
├── internal/
│   ├── handlers/        # HTTP 处理器
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑
│   ├── middleware/      # 中间件
│   └── config/          # 配置管理
└── pkg/                 # 公共包
```

### 前端结构
```
frontend/src/
├── components/          # 可复用组件
├── views/              # 页面组件
├── api/                # API 调用
├── router/             # 路由配置
└── assets/             # 静态资源
```

## 🔍 常见问题

### 开发环境问题

**Q: 启动服务时端口被占用怎么办？**
A: 使用 `./stop.sh` 停止所有服务，或者修改配置文件中的端口号。

**Q: 前端依赖安装失败？**
A: 删除 `node_modules` 和 `package-lock.json`，重新运行 `npm install`。

**Q: Go 模块下载失败？**
A: 设置 Go 代理：`go env -w GOPROXY=https://goproxy.cn,direct`

### 代码提交问题

**Q: 如何解决合并冲突？**
A: 
```bash
git fetch upstream
git rebase upstream/main
# 解决冲突后
git add .
git rebase --continue
```

**Q: 如何修改提交信息？**
A: `git commit --amend` 修改最后一次提交信息

## 📋 发布流程

项目维护者发布新版本的流程：

1. 更新版本号
2. 更新 CHANGELOG.md
3. 创建 Git 标签
4. 构建和测试
5. 发布到 GitHub Releases

## 🙏 致谢

感谢每一位贡献者！您的参与让 FlashMind 变得更好。

## 📞 联系方式

如有任何问题，请通过以下方式联系：

- **GitHub Issues**: https://github.com/lwj1989/FlashMind/issues
- **GitHub Discussions**: https://github.com/lwj1989/FlashMind/discussions
