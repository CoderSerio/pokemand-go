# pokemand-go

一个命令行工具, 用于管理脚本。

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


测试本地构建：

```bash
make build
```

测试发布流程（不会真正发布）：

```bash
make release-dry-run
```

1. 如果一切正常，继续发布：
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

