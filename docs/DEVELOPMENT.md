# 开发指南

本文档为 FlashMind 项目的开发者提供详细的开发指南。

## 🛠️ 开发环境设置

### 系统要求

- **Go**: 1.21+
- **Node.js**: 18+
- **npm**: 8+
- **Git**: 2.0+

### 本地开发设置

#### 1. 克隆仓库
```bash
git clone https://github.com/lwj1989/FlashMind.git
cd FlashMind
```

#### 2. 后端设置
```bash
cd backend
go mod download
make install-tools  # 安装开发工具
```

#### 3. 前端设置
```bash
cd frontend
npm install
```

#### 4. 启动开发环境
```bash
# 方式1: 使用脚本（推荐）
./start.sh

# 方式2: 手动启动
# 终端1 - 后端
cd backend && make run

# 终端2 - 前端
cd frontend && npm run dev
```

## 🏗️ 项目架构

### 后端架构

```
backend/
├── cmd/server/          # 应用入口点
├── internal/            # 私有代码
│   ├── config/         # 配置管理
│   ├── handlers/       # HTTP 处理器 (Controller层)
│   ├── services/       # 业务逻辑 (Service层)
│   ├── models/         # 数据模型 (Model层)
│   └── middleware/     # 中间件
└── pkg/                # 公共代码
    └── database/       # 数据库配置
```

### 前端架构

```
frontend/src/
├── api/                # API 调用
├── components/         # 可复用组件
├── views/             # 页面组件
├── router/            # 路由配置
├── services/          # 服务层
├── assets/            # 静态资源
└── test/              # 测试配置
```

## 🔧 开发工作流

### 代码规范

#### Go 代码规范
- 使用 `gofmt` 格式化代码
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `golangci-lint` 进行代码检查

```bash
# 格式化代码
make fmt

# 代码检查
make lint

# 运行测试
make test
```

#### Vue/JavaScript 代码规范
- 使用 ESLint + Prettier
- 遵循 Vue 3 Composition API 风格
- 组件名使用 PascalCase

```bash
# 格式化代码
npm run format

# 代码检查
npm run lint

# 运行测试
npm run test
```

### Git 工作流

#### 分支策略
- `main`: 主分支，保持稳定
- `develop`: 开发分支
- `feature/*`: 功能分支
- `bugfix/*`: 修复分支
- `hotfix/*`: 热修复分支

#### 提交信息规范
使用 [Conventional Commits](https://conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型：**
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或工具的变动

**示例：**
```bash
git commit -m "feat(api): add card search functionality"
git commit -m "fix(frontend): resolve display issue on mobile"
git commit -m "docs: update development guide"
```

## 🧪 测试

### 后端测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行竞态检测测试
make test-race
```

### 前端测试

```bash
# 运行单元测试
npm run test

# 运行测试并生成覆盖率报告
npm run test:coverage

# 以 UI 模式运行测试
npm run test:ui
```

### 测试编写指南

#### 后端测试
- 为每个 service 编写单元测试
- 为每个 handler 编写集成测试
- 使用 testify 库进行断言

```go
func TestCardService_CreateCard(t *testing.T) {
    // 测试代码
}
```

#### 前端测试
- 为组件编写单元测试
- 为 API 调用编写集成测试
- 使用 Vitest + Vue Test Utils

```javascript
import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import MyComponent from '@/components/MyComponent.vue'

describe('MyComponent', () => {
  it('renders properly', () => {
    // 测试代码
  })
})
```

## 🐳 Docker 开发

### 本地 Docker 开发

```bash
# 构建镜像
docker build -t flashmind .

# 运行容器
docker run -p 8080:8080 flashmind

# 使用 docker-compose
docker-compose up -d
```

### 开发时的 Docker 技巧

```bash
# 只构建特定阶段
docker build --target frontend-builder -t flashmind-frontend .

# 查看容器日志
docker-compose logs -f flashmind

# 进入容器
docker-compose exec flashmind sh
```

## 🚀 部署

### 开发环境部署

```bash
# 启动开发环境
./start.sh

# 检查状态
./status.sh

# 停止服务
./stop.sh
```

### 生产环境部署

```bash
# 构建生产版本
make prod-build

# 使用 Docker 部署
docker-compose -f docker-compose.yml up -d
```

## 🔍 调试技巧

### 后端调试

```bash
# 使用 delve 调试器
dlv debug cmd/server/main.go

# 查看日志
tail -f backend.log

# 性能分析
go tool pprof http://localhost:8080/debug/pprof/heap
```

### 前端调试

```bash
# 开启开发服务器的调试模式
npm run dev -- --debug

# 查看构建分析
npm run build -- --analyze

# 查看日志
tail -f frontend.log
```

## 📚 常用命令

### 后端常用命令
```bash
make help           # 查看所有可用命令
make build          # 构建应用
make test           # 运行测试
make lint           # 代码检查
make clean          # 清理文件
make deps           # 下载依赖
```

### 前端常用命令
```bash
npm run dev         # 开发模式
npm run build       # 构建生产版本
npm run test        # 运行测试
npm run lint        # 代码检查
npm run format      # 格式化代码
```

## 🆘 故障排除

### 常见问题

#### 端口被占用
```bash
# 查看端口占用
lsof -i :8080
lsof -i :5173

# 杀死进程
kill -9 <PID>
```

#### 依赖问题
```bash
# Go 依赖问题
go mod tidy
go clean -modcache

# Node.js 依赖问题
rm -rf node_modules package-lock.json
npm install
```

#### 数据库问题
```bash
# 重置数据库
make db-reset

# 或手动删除
rm -f backend/flashcard.db
```

### 性能优化

#### 后端性能优化
- 使用连接池
- 添加缓存层
- 数据库查询优化
- 使用 pprof 进行性能分析

#### 前端性能优化
- 代码分割
- 懒加载
- 图片优化
- Bundle 分析

## 📖 进阶主题

### 添加新的 API 端点

1. 在 `models/` 中定义数据模型
2. 在 `services/` 中实现业务逻辑
3. 在 `handlers/` 中实现 HTTP 处理器
4. 在路由中注册端点
5. 编写测试

### 添加新的前端页面

1. 在 `views/` 中创建页面组件
2. 在 `router/` 中添加路由
3. 在 `api/` 中添加 API 调用
4. 编写测试

### 数据库迁移

```bash
# 备份数据库
cp flashcard.db flashcard.db.bak

# 运行迁移（如果有）
./migrate up
```

## 🤝 贡献

请参阅 [CONTRIBUTING.md](../CONTRIBUTING.md) 了解如何为项目做出贡献。

## 📞 获取帮助

- 查看 [FAQ](../README.md#常见问题)
- 创建 [GitHub Issue](https://github.com/lwj1989/FlashMind/issues)
- 参与 [GitHub Discussions](https://github.com/lwj1989/FlashMind/discussions)
