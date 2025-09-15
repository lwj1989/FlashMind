# 更新日志

本项目的所有重要更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 新增
- 项目标准化重构
- 完善的文档体系
- 标准化的项目结构

### 修改
- 更新 README.md，提供更清晰的项目介绍
- 优化 .gitignore 文件
- 重组项目目录结构

### 删除
- 清理不必要的临时文件和构建产物
- 删除重复的文档文件

## [1.0.0] - 2024-XX-XX

### 新增
- 🧠 智能复习算法 - 基于艾宾浩斯遗忘曲线
- 📚 多格式支持 - Markdown、代码高亮、数学公式
- 🏷️ 灵活标签系统 - 多层级标签管理
- 📊 学习统计 - 详细的进度追踪和数据分析
- 🔄 批量导入导出 - 支持 TXT、CSV、JSON 格式
- 🎯 个性化学习 - 多种学习模式
- 📱 响应式设计 - 桌面和移动端支持

### 技术特性
- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue 3 + Vite + Tailwind CSS + Element Plus
- **数据库**: SQLite（可扩展到 PostgreSQL/MySQL）

### API 接口
- 卡包管理 API
- 卡片管理 API  
- 标签管理 API
- 学习会话 API
- 导入导出 API

[Unreleased]: https://github.com/lwj1989/FlashMind/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/lwj1989/FlashMind/releases/tag/v1.0.0
