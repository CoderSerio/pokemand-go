[English](./README.md)

# pkmg

`pkmg` 是一个面向人类和 coding agent 的轻量、本地优先的 skill / script 管理器。

它的目标不是做一个很重的平台，而是把已经存在的本地脚本能力收口起来：让它们更容易被复用、更容易被 agent 检索和理解，同时提供一个足够轻的 Web UI 来完成新建、编辑、复制和版本切换。

## 为什么做这个

大多数本地自动化一开始都是一堆零散脚本。

但很快就会遇到这些问题：

- 多个 skill 依赖同一类脚本
- 脚本分散在不同项目里，逐渐漂移
- agent 不知道本地已经有哪些能力可以直接复用
- 编辑、复制、回滚和版本记录都变得混乱

`pkmg` 主要解决的就是这一层本地能力治理问题：

- 把本地可复用 skill 统一放在一个位置管理
- 让 agent 能够稳定地 list / search / inspect / run
- 用一个轻量 Web UI 管理 skill
- 保持主体很轻：Go 后端 + 单页 HTML + CDN 前端依赖

## 当前功能

- 初始化本地 skill 工作区到 `data/`
- 管理 `data/scripts/` 下的脚本
- 支持 JSON 输出的列举和搜索
- 查看脚本元数据和内容预览
- 运行已管理脚本
- 启动一个轻量本地 Web UI，支持：
  - 本地 skill 列表
  - 搜索
  - 新建
  - 编辑
  - 复制
  - 版本切换
  - 打开所在目录

## 安装

### 使用 Go 安装

```bash
go install github.com/CoderSerio/pokemand-go@latest
```

### 本地开发构建

```bash
git clone https://github.com/CoderSerio/pokemand-go.git
cd pokemand-go
go build -o bin/pkmg .
```

### 本地全局软链接测试

如果你在 macOS 上，并且 `/opt/homebrew/bin` 已经在 `PATH` 里：

```bash
go build -o bin/pkmg .
ln -sfn "$(pwd)/bin/pkmg" /opt/homebrew/bin/pkmg
pkmg --version
```

如果后续要移除这个全局测试链接：

```bash
rm /opt/homebrew/bin/pkmg
```

## 快速开始

初始化工作区：

```bash
pkmg init
```

打开或创建一个 skill 脚本：

```bash
pkmg open cleanup.sh
```

列出 skill：

```bash
pkmg list
pkmg list --json
```

搜索 skill：

```bash
pkmg search cleanup
pkmg search cleanup --json
```

查看 skill：

```bash
pkmg inspect cleanup.sh
pkmg inspect cleanup.sh --json
```

运行 skill：

```bash
pkmg run cleanup.sh
pkmg run cleanup.sh arg1 arg2
```

启动 Web UI：

```bash
pkmg ui
```

## Web UI

Web UI 的设计目标就是轻量：

- 后端：Go HTTP server + WebSocket 命令通道
- 前端：单页嵌入式 HTML
- 前端依赖：尽量通过 CDN 引入

当前本地 skill 管理流程支持：

- 搜索本地 skill
- 在编辑弹窗里直接新建 skill
- 编辑已有 skill
- 按系统文件风格复制 skill
- 切回历史版本
- 打开 skill 所在目录

Skill 市场 tab 当前有意隐藏，后续再看是否重新开放。

## 面向 Agent 的使用方式

如果你想把 `pkmg` 接到 agent 工作流里，优先建议走结构化 CLI，而不是先让 agent 去解析 UI。

推荐命令：

```bash
pkmg list --json
pkmg search "<query>" --json
pkmg inspect "<relative-path>" --json
pkmg run "<relative-path>" [args...]
```

这样 agent 可以稳定地拿到本地能力清单和元数据。

如果你要给仓库内 agent 一份更明确的说明，可以看 [AGENTS.md](./AGENTS.md)。

## 项目结构

```text
cmd/              Cobra 命令和后端逻辑
cmd/webui/        内嵌 Web UI 资源
data/scripts/     已管理的本地 skill 脚本
data/.pkmg/       本地元数据和版本快照
platform/         预留给未来平台分发包装层
```

## 测试

本地构建和冒烟测试：

```bash
go build ./...
pkmg init
pkmg list --json
pkmg ui
```

比较有价值的手工检查：

- 在 UI 里新建一个 skill
- 编辑并保存
- 复制 skill
- 切回历史版本
- 确认 `inspect --json` 反映的是最新状态

## 分发思路

核心产品建议始终保持为 Go 二进制。

推荐的分发顺序：

1. `go install`
2. GitHub Releases 多平台二进制
3. Homebrew
4. Windows 包管理器
5. 可选的 npm 薄包装

如果以后接 npm，建议只把它当成 Go 二进制的分发包装层，而不是用 JS 重写核心逻辑。

## 开发

构建：

```bash
go build ./...
```

本地运行：

```bash
go run . --help
go run . ui
```

格式化：

```bash
gofmt -w cmd/*.go main.go
```

## 当前阶段

`pkmg` 现在还在早期阶段，重点在：

- 把本地 skill 管理做好
- 给 agent 稳定的 discovery / inspect 能力
- 保持整体足够轻，而不是过早平台化

欢迎继续迭代和收敛方向。
