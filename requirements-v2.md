# 📚 Flashcard 项目需求文档（卡包版 v2）

## 1. 项目概述

- 名称：Flashcard（候选：GoFlashcard / CodeFlashcard）
- 目标：面向程序员的本地离线知识记忆应用，提供高效的间隔重复学习体验与优雅的使用界面。
- 价值主张：
  - 本地私有、轻量快速、零依赖云服务
  - 以卡包为核心的清晰知识组织（主题 → 分类 → 卡片）
  - 科学的 SM-2 复习调度
  - 易用的 TXT 批量导入导出，便于迁移与备份
- 运行模式：
  - 前端：Vue 3 + Vite + TailwindCSS（建议配合 Naive UI / Element Plus）
  - 后端：Go + Gin + GORM + SQLite（单文件 DB）
  - 数据：SQLite（主存），TXT（备份/迁移）
- 非功能性目标（NFR）：
  - 离线优先、本地单用户、数据持久化
  - 1 万张卡片内应保持流畅（查询与抽卡 < 200ms）
  - 良好键盘可用性（学习页尽量全键盘操作）
  - 重要操作可撤销（删除/导入回滚）

## 2. 核心概念与范围

- 卡包（Deck）：围绕一个知识主题的集合，例如「Go语言」「系统设计」。
- 标签（Tag）：卡包内的细分类别，例如「并发」「数据库」「网络」。
- 卡片（Card）：包含问题、Markdown 答案、所属卡包与标签、复习调度参数。
- 复习调度（SRS）：采用 SM-2 算法（记得/模糊/忘记三档反馈），存储 efactor、interval、repetitions、next_review。
- 范围说明：
  - 单机单用户（不涉及云同步、多端冲突）
  - 先不做多媒体（音频/图片）与 Cloze 完形等高级题型（可留待后续）

## 3. 功能需求

### 3.1 卡包（Deck）管理

- 新建、重命名、删除卡包
- 查看卡包统计：总卡片数、待复习卡片数、标签数量
- 导入到指定卡包、按卡包导出
- 可选：归档卡包（只读，不参与抽卡）

验收要点：
- 新建后主页可见；删除时要求二次确认；重命名实时生效
- 统计信息在主页卡片上展示，进入详情页可查看细分统计
- 归档卡包在学习入口隐藏，但数据可导出

### 3.2 标签（Tag）管理

- 在卡包内创建、重命名、删除标签
- 标签级学习（仅抽取该标签的待复习卡片）
- 可选：跨卡包复用同名标签（默认不跨包）

验收要点：
- 标签名在同一卡包中唯一
- 删除标签时可选择「保留卡片并置为未分组」或「连同卡片一起删除」

### 3.3 卡片（Card）管理

- 结构：问题（纯文本/Markdown）、答案（Markdown）、所属卡包、所属标签、复习参数
- CRUD：创建、编辑、删除、查看；支持批量删除与标签更换
- 搜索：按问题/答案关键词、卡包、标签筛选
- Markdown 渲染与代码高亮（`markdown-it` + `highlight.js`）

验收要点：
- 必填项校验（问题/答案非空、卡包必选）
- 编辑保存后内容与渲染一致
- 搜索速度：1 万卡内 < 200ms（本地模糊查询 + 前端缓存）

### 3.4 学习模式（SRS）

- 抽卡：按卡包或标签抽取「今日待复习」卡片
- 交互：正面显示问题 → 翻转显示答案（CSS 3D 动画）
- 反馈按钮：记得（Good）/ 模糊（Hard）/ 忘记（Again）
- 进度：显示已复习/剩余数量，今日学习数据
- 键盘操作：如 空格翻转、1=Again、2=Hard、3=Good、H/J/K/L 导航（可配置）
- 可选：每日新卡上限、学习次序（到期优先/混洗）

验收要点：
- 抽取逻辑正确：仅抽取 `next_review <= 今日` 的卡片
- 反馈后即时更新复习参数与 `next_review`
- 学习中断后可恢复至上一张卡（轻量会话持久化）

### 3.5 导入与导出（TXT）

- 文件编码：UTF-8
- 语法规则：
  - 导入时选择目标卡包（TXT 中不写卡包名）
  - 标签用 `#标签名`
  - 问题/答案用 `---` 分隔，卡片之间用 `===` 分隔
  - 问题/答案均支持 Markdown
- 示例：

```markdown
#并发
什么是 goroutine？
---
goroutine 是 Go 语言中轻量级的线程，由 Go runtime 调度。
===
Go 中的 channel 有什么作用？
---
channel 是 goroutine 之间通信的机制，用于数据传递和同步。
```

- 导入行为：
  - 未存在的标签自动创建到该卡包
  - 重复检测：同一卡包内问题+答案完全一致视为重复，默认跳过（提供统计与报告）
  - 失败回滚：导入失败需整体回滚（事务）
- 导出行为：
  - 按卡包导出，按标签分组；保留 Markdown
  - 包含导出报告（导出时间、卡包名、标签列表、卡片数）

验收要点：
- 大文件（例如 1 万卡）导入时间可接受（分批事务 + 流式解析）
- 提供导入预览（解析摘要：标签数、卡片数、重复数）

### 3.6 记忆算法（SM-2）

- 采用 SM-2（可选使用 [github.com/matthayes/sm2](https://pkg.go.dev/github.com/matthayes/sm2)）
- 存储字段：`efactor`（浮点）、`interval`（天）、`repetitions`（连续复习次数）、`next_review`（时间）
- 反馈规则：
  - Again：重置 `repetitions`，`interval` 设为 1（或更小），`efactor` 下调
  - Hard：小幅上调 `interval`，`efactor` 微调
  - Good：正常递增 `interval`，`efactor` 小幅提升/保持
- 默认初值：`efactor=2.5`，`interval=0`，`repetitions=0`
- 时区：以本地时区计算「今日」

验收要点：
- 连续多次 Hard 不应导致 `efactor` 过低（设下限例如 1.3）
- next_review 必为日期起始（避免跨日边界误差）

### 3.7 统计与报告

- 主页：每个卡包显示总卡、待复习卡、标签数
- 详情页：按标签/今日/本周完成数、学习时长（可选）
- 导入导出报告：成功/跳过/失败条目统计

## 4. 系统架构与接口

### 4.1 后端（Go + Gin + GORM + SQLite）

- 主要模块：
  - Decks：卡包管理
  - Tags：标签管理
  - Cards：卡片 CRUD
  - Study：学习调度（取队列、提交反馈）
  - Import/Export：TXT 解析、导入导出
  - Stats：统计
- API（示例）：
  - GET `/decks`：列出卡包
  - POST `/decks`：创建卡包
  - PATCH `/decks/:deckID`：重命名/归档
  - DELETE `/decks/:deckID`：删除卡包
  - GET `/decks/:deckID/tags`：列出标签
  - POST `/decks/:deckID/tags`：创建标签
  - PATCH `/tags/:tagID`：重命名
  - DELETE `/tags/:tagID`：删除
  - GET `/decks/:deckID/cards`：分页查询（支持搜索/筛选）
  - POST `/decks/:deckID/cards`：创建卡片
  - PATCH `/cards/:id`：更新卡片
  - DELETE `/cards/:id`：删除卡片
  - GET `/study/queue?deckID=&tagID=`：获取今日待复习队列
  - POST `/cards/:id/review`：提交复习结果（Again/Hard/Good）
  - POST `/decks/:deckID/import`：导入 TXT
  - GET `/decks/:deckID/export`：导出 TXT
  - GET `/decks/:deckID/stats`：卡包统计
- 错误返回：
  - 统一结构 `{ code, message, details? }`
  - 常见 code：`InvalidParam`、`NotFound`、`Conflict`、`Internal`

### 4.2 数据库设计（SQLite）

- 表 `decks`：`id` PK，`name` UNIQUE，`archived` BOOL，`created_at`
- 表 `tags`：`id` PK，`deck_id` FK->decks，`name`（同包内唯一），`created_at`
- 表 `cards`：`id` PK，`deck_id` FK，`tag_id` FK（可空=未分组），`question`，`answer`，`created_at`，`updated_at`
- 表 `reviews`：`id` PK，`card_id` FK UNIQUE（1:1 当前调度状态），`efactor` FLOAT，`interval` INT，`repetitions` INT，`next_review` DATETIME
- 约束与索引：
  - 索引：`cards(deck_id, tag_id)`，`reviews(next_review)`，`cards(question)`（可选 FTS）
  - 约束：`ON DELETE CASCADE`（删除卡包/标签时联动）
  - 去重：同一 `deck_id` 下 `question+answer` 组合唯一（导入时可软校验）

## 5. 前端设计（Vue 3 + Vite + TailwindCSS）

- 页面与功能：
  - 首页（卡包总览）：卡片式卡包，展示统计，入口按钮（学习/管理/导入导出）
  - 卡包详情页：标签列表与统计、卡片列表（搜索/分页/批量操作）
  - 标签详情页：同卡包视图但限定标签筛选
  - 学习页：问题→翻转→答案，三档反馈，进度条，键盘操作
  - 导入导出页：TXT 上传、语法提示、预览与结果报告、导出按钮
  - 设置页（可选）：主题切换、快捷键映射
- UI 建议：
  - Notion/Anki 风格，清爽、支持深色
  - 组件库：Naive UI / Element Plus
  - 动效：CSS 3D 翻转
  - Markdown：`markdown-it` + `highlight.js`

## 6. 非功能性与工程约束

- 性能：1 万卡以内主要操作 < 200ms；导入 1 万卡 < 数分钟
- 稳定性：导入使用事务；关键操作支持撤销或提供导出前快照
- 兼容：Chrome/Edge/Safari 最新版；macOS/Windows
- 安全：默认仅本地访问；可配置跨域关闭；无第三方上报
- 日志：后端结构化日志（level、trace id）；导入导出生成报告
- 配置：`.env`（数据库文件路径、端口、导入导出目录等）

## 7. 备份与恢复

- 一键导出卡包为 TXT（含标签分组）
- 导入 TXT 即可恢复/迁移
- 可选：后端 CLI 定时导出（crontab）

## 8. 验收标准（关键场景）

- 按卡包/标签抽卡仅出现今日应复习卡片，反馈后参数及下次复习时间正确更新
- TXT 导入严格按语法解析，事务回滚可用，重复项处理符合策略
- 主页与详情统计准确、实时
- 学习页键盘操作可替代鼠标，翻转动画流畅
- 暗色模式、Markdown 渲染和代码高亮表现正确

## 9. 开放问题（待决策）

- 是否支持「Easy」第四档反馈？（当前为 Again/Hard/Good）
- 标签是否允许跨卡包复用并统一统计？（默认否）
- 删除标签时默认策略是仅移除标签还是连同卡片删除？
- 搜索是否需要全文检索（FTS5）还是普通 LIKE 即可？
- 是否需要「撤销上一次复习」功能（以及其限制）？

---

# 🧩 需求拆分（可执行工作分解）

## 阶段 1：基础设施与数据库

- 建仓与工程脚手架（前后端目录、Docker 可选、Makefile/Taskfile）
- 数据库建模与迁移（decks/tags/cards/reviews、索引、约束）
- 通用错误结构、中间件（日志、恢复、CORS/本地）

交付物：可启动的后端服务，空白前端应用，数据库初始化成功

## 阶段 2：卡包与标签

- 后端：`/decks` CRUD、`/decks/:deckID/tags` CRUD、统计接口
- 前端：主页卡包列表卡片、卡包详情页（标签列表+统计）
- 约束：同包标签名唯一；删除标签的处理策略

交付物：能管理卡包与标签，主页可见统计

## 阶段 3：卡片 CRUD 与搜索

- 后端：`/decks/:deckID/cards` 列表与创建、`/cards/:id` 更新/删除、搜索参数
- 前端：卡片列表、编辑弹窗、批量操作、Markdown 渲染与高亮
- 性能：分页与前端缓存策略

交付物：卡片可增删改查并可用搜索

## 阶段 4：导入导出（TXT）

- 后端：解析器、导入事务、重复检测、导入预览、导出
- 前端：导入向导（选择卡包→上传→预览→执行→报告）、导出按钮
- 可选：CLI 定时导出

交付物：稳定的导入导出闭环，失败回滚可靠

## 阶段 5：学习页与 SM-2

- 后端：`/study/queue`、`/cards/:id/review`（实现 SM-2）
- 前端：翻转交互、三档反馈、进度统计、键盘操作、每日上限
- 边界：本地时区处理、EF 下限、撤销上次复习（可选）

交付物：可连续学习并正确调度的学习体验

## 阶段 6：UI 完善与体验优化

- 深色模式、CSS 3D 动画打磨、空态与加载态
- 错误提示与表单校验文案统一
- 统计可视化（简单图表可选）

交付物：一致、顺滑的用户体验

## 阶段 7：测试与发布

- 后端：单元测试（导入解析、SM-2 计算）、接口集成测试
- 前端：关键交互测试（学习页、导入流程）
- 文档：使用说明、导入语法、备份恢复指南

交付物：测试通过与可分发包（或运行指南）

---

## 建议的任务清单（优先级）

- 初始化后端与数据库迁移脚手架（高）
- 实现 `decks/tags` 模块与统计（高）
- 实现 `cards` CRUD 与搜索（高）
- 实现 TXT 导入解析与事务化导入（高）
- 实现 TXT 导出（高）
- 实现 `study/queue` 与 `review`（SM-2）（高）
- 学习页前端交互（翻转与键盘）（中）
- 深色模式与 UI 打磨（中）
- CLI 定时导出（可选）（低）
- 统计可视化（可选）（低）


