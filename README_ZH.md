# <div align="center">pkmg</div>

<div align="center">

_给人类和 coding agent 用的，本地优先 skill 管理器。_

[![release](https://img.shields.io/github/v/release/CoderSerio/pokemand-go?display_name=tag&style=flat-square)](https://github.com/CoderSerio/pokemand-go/releases)
[![ci](https://img.shields.io/github/actions/workflow/status/CoderSerio/pokemand-go/ci.yml?branch=main&style=flat-square&label=ci)](https://github.com/CoderSerio/pokemand-go/actions/workflows/ci.yml)
[![go](https://img.shields.io/github/go-mod/go-version/CoderSerio/pokemand-go?style=flat-square)](https://go.dev/)
[![license](https://img.shields.io/github/license/CoderSerio/pokemand-go?style=flat-square)](./LICENSE)
[![local-first](https://img.shields.io/badge/local--first-yes-0f766e?style=flat-square)](#为什么是-pkmg)
[![agent-ready](https://img.shields.io/badge/agent-ready-yes-1d4ed8?style=flat-square)](./AGENTS.md)

[English](./README.md) · [快速开始](#快速开始) · [Web UI](#web-ui) · [Agent 用法](#面向-agent-的使用方式) · [测试](#测试)

</div>

`pkmg` 把四散在本地的脚本能力收成一套可复用的 skill 库。

它更像是本地 agent 工作流里的轻量脚本层管理器：负责检索、查看、编辑、版本化和运行已有 skill，而不是去做一个重平台。

`pkmg` 是 `pokemand-go` 的命令行入口名，可以理解成一个用 Go 写的“口袋命令”管理器。

## 为什么是 pkmg

现在 `skill` 确实越来越像 agent 时代更自然的本地化方案。

但本地化之后，还是会留下一个很具体的问题：

- skill 里的脚本会越来越多
- 同类能力会在不同项目和 skill 之间反复复制
- 现有本地能力越来越难查、越来越难管
- agent 想复用这些能力时，缺的是结构化检索入口，而不是更多目录

`pkmg` 处理的就是这一层。

它不替代 skill，而是为 skill 服务：把脚本这层能力统一管理起来，一次封装，到处复用，并且让 agent 能轻松上手已有能力。

## 它现在是什么感觉

- 默认本地优先，不依赖云端注册表
- 从一开始就对 agent 友好，CLI 支持结构化 JSON 输出
- 主体很轻，Go 后端 + 嵌入式页面 + CDN 前端依赖
- 自带版本意识，脚本保存后会留下本地快照
- 人也能顺手用，一个小 Web UI 就能完成新建、编辑、复制、回滚、打开目录

## 当前能力

- 初始化用户级 skill 工作区
- 在独立的数据目录下管理脚本
- 通过 CLI 完成 list / search / inspect / run
- 为 agent 工作流提供结构化 JSON 输出
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

### 从源码构建

```bash
git clone https://github.com/CoderSerio/pokemand-go.git
cd pokemand-go
go build -o bin/pkmg .
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

## 默认目录模型

`pkmg` 现在默认使用用户级目录，而不是仓库内的 `data/`。

- 配置根目录：`os.UserConfigDir()/pkmg`
- 数据根目录：`PKMG_DATA_DIR`，或配置中的 `dataPath`，否则回退到 `os.UserConfigDir()/pkmg`
- 脚本目录：`<data-root>/scripts`
- 版本快照：`<data-root>/.pkmg`

如果你想改到自己的位置：

```bash
export PKMG_CONFIG_DIR=/your/custom/config
export PKMG_DATA_DIR=/your/custom/data
```

这样默认就不会把托管 skill 数据塞进仓库本身，更适合做可复用的本地工具层。

## Web UI

- 后端：Go HTTP server + WebSocket 命令通道
- 前端：单页嵌入式页面
- 前端依赖：优先通过 CDN 引入，避免把主体做重

当前本地 skill 流程支持：

- 搜索本地 skill
- 在编辑弹窗里直接新建 skill
- 用轻量代码视图编辑已有 skill
- 按系统文件风格复制 skill
- 切回历史版本
- 在系统文件浏览器中打开所在目录

## 面向 Agent 的使用方式

如果你要把 `pkmg` 接进 agent 工作流，优先建议走结构化 CLI，而不是让 agent 去解析 UI。

推荐命令：

```bash
pkmg list --json
pkmg search "<query>" --json
pkmg inspect "<relative-path>" --json
pkmg run "<relative-path>" [args...]
```

这样 agent 可以稳定获得本地能力清单，而不是猜目录结构。

更具体的 agent 集成说明见 [AGENTS.md](./AGENTS.md)。

## 测试

自动化：

```bash
go test ./...
```

快速冒烟：

```bash
go build ./...
pkmg init
pkmg list --json
pkmg ui
```
