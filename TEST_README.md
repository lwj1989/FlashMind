# FlashMind 测试指南

## 测试数据

我们提供了Golang基础知识和面试题的测试数据，位于 `test_data/go_flashcards.json` 文件中。该数据包含3个卡包：

1. **Go语言基础** - 包含Go语言基础知识的卡片
2. **Go语言面试题** - 包含常见的Go语言面试题
3. **Go语言高级特性** - 包含Go语言高级特性的卡片

每个卡包包含多个卡片，每张卡片有问题、答案和标签信息。

## 运行测试

### 前端测试

前端测试使用Vitest框架，包括组件测试和API测试。

#### 安装依赖
```bash
cd frontend
npm install
```

#### 运行所有测试
```bash
npm run test
```

#### 运行测试并生成覆盖率报告
```bash
npm run test:coverage
```

#### 运行测试UI界面
```bash
npm run test:ui
```

### 后端测试

后端测试使用Go的内置测试框架。

#### 运行所有测试
```bash
cd backend
go test -v ./...
```

#### 运行测试并生成覆盖率报告
```bash
cd backend
go test -v -cover ./...
```

### 运行所有测试

使用提供的脚本可以一次性运行所有测试：

```bash
./run_tests.sh
```

## 测试用例说明

### 前端测试用例

#### 组件测试
- `CardForm.spec.js` - 测试卡片表单组件
- `CardItem.spec.js` - 测试卡片项组件
- `CardList.spec.js` - 测试卡片列表组件

#### 页面测试
- `HomeView.spec.js` - 测试首页
- `CardsView.spec.js` - 测试卡片管理页面
- `CardDetailView.spec.js` - 测试卡片详情页面
- `DecksView.spec.js` - 测试卡包管理页面
- `DeckDetailView.spec.js` - 测试卡包详情页面
- `TagsView.spec.js` - 测试标签管理页面
- `TagCardsView.spec.js` - 测试标签卡片页面
- `ImportExportView.spec.js` - 测试导入导出页面

#### API测试
- `request.spec.js` - 测试请求配置
- `deck.spec.js` - 测试卡包API
- `card.spec.js` - 测试卡片API
- `tag.spec.js` - 测试标签API
- `importExport.spec.js` - 测试导入导出API

### 后端测试用例

#### 处理器测试
- `deck_test.go` - 测试卡包处理器
- `card_test.go` - 测试卡片处理器
- `tag_test.go` - 测试标签处理器
- `import_export_test.go` - 测试导入导出处理器

## 测试覆盖率

测试覆盖率报告将生成在以下位置：

- 前端：`frontend/coverage/`
- 后端：`backend/coverage/`

## 使用测试数据导入

1. 启动后端服务器
2. 启动前端服务器
3. 访问导入导出页面
4. 使用 `test_data/go_flashcards.json` 文件导入测试数据

## 注意事项

- 确保在运行测试前已安装所有依赖
- 运行测试前确保没有其他进程占用测试端口
- 如果测试失败，请检查错误日志并确保数据库配置正确