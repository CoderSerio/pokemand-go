# Pokemand-go

一个命令行工具。

## 安装方法

### 方法1：通过 go install 安装（推荐）

确保你已经安装了 Go 1.24 或更高版本，然后运行：

```bash
go install github.com/CoderSerio/pokemand-go@latest
```

安装后会自动创建 `pkmg` 命令。

### 方法2：下载二进制文件

1. 从 [Releases](https://github.com/CoderSerio/pokemand-go/releases) 页面下载对应你操作系统的二进制文件
2. 解压下载的文件
3. 将二进制文件移动到系统 PATH 目录下：

Linux/macOS:
```bash
chmod +x pkmg-darwin-amd64
sudo mv pkmg-darwin-amd64 /usr/local/bin/pkmg
```

Windows:
- 将 `pkmg-windows-amd64.exe` 重命名为 `pkmg.exe`
- 移动到 `C:\Windows\System32` 或其他在 PATH 中的目录

### 初始化

安装完成后，需要初始化环境：

```bash
pkmg init
```

可选：创建命令别名（推荐）：

```bash
pkmg init --aliases
```

## 使用方法

```bash
# 查看帮助
pkmg --help

# 查看版本
pkmg version
```

## 开发

### 构建

```bash
# 本地构建
make build

# 安装到 GOPATH/bin
make install

# 测试发布流程（不会真正发布）
make release-dry-run

# 发布新版本（需要先创建 git tag）
make release
```

### 发布新版本

```bash
# 1. 创建新的版本标签
git tag -a v0.1.0 -m "First release"

# 2. 推送标签到 GitHub，这会自动触发发布流程
git push origin v0.1.0
```
```

主要变更说明：

1. **GoReleaser 配置**：
   - 设置了与 Makefile 相同的二进制名称 (pkmg)
   - 使用相同的 ldflags 注入版本信息
   - 添加了发布说明模板
   - 配置了与你的 GitHub 仓库匹配的所有者和仓库名

2. **Makefile 简化**：
   - 移除了重复的发布逻辑
   - 添加了使用 GoReleaser 的命令
   - 保留了本地开发需要的命令

3. **工作流程**：
   - 本地开发时使用 `make build` 和 `make install`
   - 要发布新版本时，创建并推送 tag 即可
   - GitHub Actions 会自动处理发布流程

要测试这些更改：

1. 首先测试本地构建：
```bash
make build
```

2. 测试发布流程（不会真正发布）：
```bash
make release-dry-run
```

3. 如果一切正常，可以创建第一个发布：
```bash
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

这样配置后，发布流程会：
1. 自动构建所有平台的二进制文件
2. 创建 checksums 文件
3. 生成 changelog
4. 创建 GitHub release
5. 上传所有构建的文件

需要我解释任何部分吗？或者你想先测试看看效果？

### 构建

```bash
# 本地构建
make build

# 安装到 GOPATH/bin
make install

# 构建多平台发布包
make release
``` 