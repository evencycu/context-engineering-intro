# Cursor Rules 使用指南

## 概述

本项目已配置增强版的 `.cursorrules` 文件，实现了稳定的增量开发流程。Cursor AI 助手会按照这些规则自动执行测试、更新文档等任务。

## 核心功能

### 1. 增量开发工作流程

当开发新功能或重构代码时，Cursor 会自动执行以下步骤：

1. **规划阶段** - 理解需求、检查现有实现
2. **开发阶段** - 编写代码、添加注释、错误处理
3. **测试阶段** - 编写并运行测试、检查覆盖率
4. **文档阶段** - 更新相关文档、验证同步
5. **代码质量** - 格式化、静态分析、lint 检查
6. **提交阶段** - 使用规范的提交信息

### 2. 自动化检查

- ✅ 代码格式化 (`go fmt`)
- ✅ 静态分析 (`go vet`)
- ✅ 代码检查 (`golangci-lint`)
- ✅ 测试执行 (`go test`)
- ✅ 覆盖率检查 (≥80%)
- ✅ 文档同步验证

## 使用方法

### 方式 1: 使用 Makefile 命令

```bash
# 在 Backend 目录下执行

# 提交前检查
make dev-pre-commit

# 开发新功能时
make dev-feature

# 重构代码时
make dev-refactor

# 完整工作流程
make dev-full

# 单独检查代码质量
make dev-quality

# 单独检查文档同步
make dev-docs
```

### 方式 2: 直接使用脚本

```bash
# 在项目根目录执行

# 提交前检查
./scripts/dev-workflow.sh pre-commit

# 开发新功能时
./scripts/dev-workflow.sh feature

# 重构代码时
./scripts/dev-workflow.sh refactor

# 完整工作流程
./scripts/dev-workflow.sh full

# 单独运行测试
./scripts/dev-workflow.sh test

# 单独检查代码质量
./scripts/dev-workflow.sh quality

# 单独检查文档
./scripts/dev-workflow.sh docs
```

### 方式 3: 让 Cursor AI 自动执行

直接告诉 Cursor：

- "帮我开发一个新功能：..."
- "帮我重构这段代码：..."
- "帮我修复这个 bug：..."

Cursor 会自动按照 `.cursorrules` 中的工作流程执行。

## 工作流程详解

### 功能开发流程

```
1. 规划阶段
   ├─ 理解需求和范围
   ├─ 检查现有类似实现
   ├─ 规划 API 变更（如有）
   └─ 识别受影响文件和依赖

2. 开发阶段
   ├─ 编写代码（遵循项目风格）
   ├─ 添加注释（复杂逻辑）
   ├─ 错误处理
   └─ 原子性提交

3. 测试阶段 ⚠️ 必需
   ├─ 编写单元测试
   ├─ 编写集成测试
   ├─ 运行所有测试
   ├─ 检查覆盖率（≥80%）
   ├─ 测试边界情况
   └─ E2E 测试验证

4. 文档阶段 ⚠️ 必需
   ├─ 更新 API 文档
   ├─ 更新用户指南
   ├─ 更新架构文档
   ├─ 运行文档同步检查
   └─ 添加代码示例

5. 代码质量阶段
   ├─ 格式化代码
   ├─ 运行静态分析
   ├─ 运行 linter
   └─ 代码审查

6. 提交阶段
   └─ 使用规范格式提交
```

### 重构流程

```
1. 确保现有测试通过
2. 添加新测试（针对重构代码）
3. 保持向后兼容性
4. 更新受影响文档
5. 运行完整测试套件
6. 验证无性能回归
```

### Bug 修复流程

```
1. 编写测试重现 bug
2. 修复 bug
3. 确保测试通过
4. 检查类似问题
5. 更新文档（如行为改变）
```

## 测试要求

### 覆盖率标准

- **单元测试**: ≥80% 新代码覆盖率
- **集成测试**: 所有 API 端点必需
- **E2E 测试**: 关键用户流程必需
- **负载测试**: 性能关键功能必需

### 测试最佳实践

- ✅ 测试优先或并行编写（TDD 优先）
- ✅ 测试成功和失败场景
- ✅ 使用描述性测试名称
- ✅ 保持测试独立和隔离
- ✅ 使用表驱动测试
- ✅ Mock 外部依赖

## 文档要求

### 何时更新文档

- ✅ 添加/修改 API 端点时
- ✅ 更改用户面向功能时
- ✅ 修改架构或设计模式时
- ✅ 添加新配置选项时
- ✅ 更改部署流程时

### 文档结构

```
docs/
├── 01_REQUIREMENTS/    # 需求文档
├── 02_ARCHITECTURE/    # 架构文档
├── 03_DESIGN/          # 设计文档（API 规范）
├── 04_TEST/            # 测试文档
├── 05_DEPLOYMENT/      # 部署文档
├── 06_USER_GUIDE/      # 用户指南
└── 07_DEVELOPMENT/     # 开发指南
```

## 代码规范

### Go 代码

- 使用 `golangci-lint` 检查代码质量
- 遵循 Go 命名约定
- 使用 `go fmt` 格式化
- 显式错误处理
- 使用 `context.Context` 处理取消和超时

### React/TypeScript 代码

- 使用 TypeScript 严格模式
- 优先使用函数组件和 hooks
- 使用 ES 模块语法
- 尽可能解构导入
- 遵循 React 最佳实践

## Git 提交规范

### 格式

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

### 示例

```
feat(api): add external notification endpoint

- Add POST /api/v1/notify endpoint
- Support multiple notification types
- Add rate limiting

Closes #123
```

## 快速参考

### 常用命令

```bash
# 启动开发环境
make docker-up && make db-init && make run

# 运行测试
make test

# 检查代码质量
make dev-quality

# 更新文档
make dev-docs

# 完整开发周期
make dev-full
```

### 开发检查清单

开发新功能时，确保：

- [ ] 代码已格式化 (`go fmt`)
- [ ] 所有测试通过 (`make test`)
- [ ] 覆盖率 ≥80%
- [ ] 文档已更新
- [ ] 文档同步检查通过
- [ ] 提交信息符合规范

## 故障排除

### 测试失败

```bash
# 运行详细测试输出
go test -v ./...

# 运行特定包的测试
go test ./services/teamsnotification/...

# 检查竞态条件
go test -race ./...
```

### 文档同步问题

```bash
# 手动运行文档同步检查
./scripts/utils/check-doc-sync.sh

# 查看检查报告
cat doc_sync_report.md
```

### Linter 错误

```bash
# 运行 linter
golangci-lint run

# 自动修复可修复的问题
golangci-lint run --fix
```

## 最佳实践

1. **测试优先**: 编写测试优先或并行编写代码
2. **文档始终**: 每次变更都更新文档
3. **增量开发**: 小而专注的变更，伴随测试
4. **质量保证**: 提交前格式化、lint 和测试
5. **一致性**: 遵循项目约定和模式

## 相关文件

- `.cursorrules` - Cursor AI 行为规则
- `scripts/dev-workflow.sh` - 开发工作流程脚本
- `Backend/Makefile` - Makefile 命令
- `Backend/docs/07_DEVELOPMENT/` - 开发文档

## 参考资源

- [Cursor Rules Documentation](https://docs.cursor.com/context/rules)
- [Claude.md Best Practices](https://arxiv.org/abs/2509.14744)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [React Best Practices](https://react.dev/learn)
