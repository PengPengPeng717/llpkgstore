# llpkgstore 项目文档

[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/llpkgstore)](https://goreportcard.com/report/github.com/goplus/llpkgstore)
[![Go Version](https://img.shields.io/github/go-mod/go-version/goplus/llpkgstore)](https://github.com/goplus/llpkgstore/blob/main/go.mod)
[![License](https://img.shields.io/github/license/goplus/llpkgstore)](https://github.com/goplus/llpkgstore/blob/main/LICENSE)

## 📖 目录

- [项目概述](#项目概述)
- [核心特性](#核心特性)
- [系统架构](#系统架构)
- [安装与部署](#安装与部署)
- [快速开始](#快速开始)
- [详细使用指南](#详细使用指南)
- [配置说明](#配置说明)
- [命令参考](#命令参考)
- [常见问题](#常见问题)
- [故障排除](#故障排除)
- [贡献指南](#贡献指南)
- [更新日志](#更新日志)

## 🎯 项目概述

**llpkgstore** 是一个专为 [**LLGo**](https://github.com/goplus/llgo) 设计的综合包分发服务，旨在为多语言生态系统提供可信赖且便捷的语言绑定访问。

### 项目背景

随着跨语言开发的普及，开发者需要在不同编程语言之间进行互操作。LLGo 作为 Go 语言的扩展，提供了强大的跨语言集成能力。llpkgstore 在此基础上，为开发者提供了一个统一的工具链，用于生成、管理和分发各种语言的 Go 绑定。

### 设计理念

- **统一性**: 提供一致的接口和体验
- **可信赖性**: 通过自动化流程确保包质量
- **易用性**: 简化复杂的跨语言绑定生成过程
- **可扩展性**: 支持多种编程语言和包管理器

## ✨ 核心特性

### 🌍 多语言支持

| 语言/平台 | 支持状态 | 生成工具 | 包管理器 | 特性 |
|-----------|----------|----------|----------|------|
| **C/C++** | ✅ 完全支持 | `llcppg` | Conan | 二进制分发、头文件、.pc 文件 |
| **Python** | ✅ 完全支持 | `llpyg` | pip | 模块绑定、Go 接口、类型安全 |
| **JavaScript** | 🚧 计划中 | - | - | 计划支持 |

### 🔧 核心能力

- **自动化生成**: 完整的 CI/CD 流水线，自动创建包
- **统一版本管理**: C/C++ 和 Python 包使用一致的版本映射机制
- **智能版本提取**: 从提交消息自动提取版本信息，支持多种格式
- **双重版本记录**: 同时维护本地和集中式版本记录文件
- **自动 Git 标签**: 基于版本信息自动创建和推送 Git 标签
- **类型安全**: 跨语言边界保持类型信息
- **测试框架**: 自动生成演示代码和测试
- **GitHub 集成**: 无缝的 GitHub Actions 和 Releases 支持
- **智能路由**: 根据包类型自动选择相应的处理逻辑

### 🎯 应用场景

- **库集成**: 在 Go 中使用 C/C++ 库
- **Python 生态**: 从 Go 访问 Python 包
- **跨平台**: 在不同平台上提供一致的体验
- **企业级**: 生产级的包管理解决方案

## 🏗️ 系统架构

### 整体架构图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户界面层     │    │   业务逻辑层     │    │   数据存储层     │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│  CLI 命令接口   │    │  包类型检测      │    │  配置文件       │
│  Web 服务接口   │    │  命令路由        │    │  Git 仓库       │
│  API 接口       │    │  包生成逻辑      │    │  元数据存储     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   外部工具集成   │
                    ├─────────────────┤
                    │  llcppg (C/C++) │
                    │  llpyg (Python) │
                    │  Conan/pip      │
                    └─────────────────┘
```

### 核心组件

#### 1. 主程序入口 (`main.go`)
- 负责包类型检测和命令路由
- 根据 `llpkg.cfg` 中的 `type` 字段选择处理逻辑
- 支持智能回退机制

#### 2. 命令实现层
- **C/C++ 包**: `internal_cpp/` 目录下的命令实现
- **Python 包**: `internal_python/` 目录下的命令实现

#### 3. 配置管理
- 统一的配置文件格式 (`llpkg.cfg`)
- 支持包类型特定的配置选项
- 自动配置验证和错误提示

#### 4. 版本管理
- **统一版本映射**: C/C++ 和 Python 包使用一致的版本管理机制
- **智能版本提取**: 从提交消息自动提取版本信息
- **双重版本记录**: 同时更新本地和集中式版本记录文件
- **自动 Git 标签**: 基于版本信息自动创建和推送 Git 标签
- **版本验证**: 完整的版本格式验证和冲突检测

## 🚀 安装与部署

### 系统要求

- **Go**: 1.24.x 或更高版本
- **Python**: 3.12+ (用于 Python 包支持)
- **Git**: 最新版本
- **操作系统**: Linux, macOS

### 依赖工具

#### 必需工具
- **LLGo**: 用于 Python 绑定生成
- **LLPyg**: Python 绑定生成工具
- **Conan**: C/C++ 包管理 (C++ 包)

#### 可选工具
- **GitHub CLI**: 用于 GitHub 操作
- **Docker**: 容器化部署

### 安装方法

#### 方法 1: 从源码构建

```bash
# 克隆仓库
git clone https://github.com/PengPengPeng717/llpkgstore.git
# git clone https://github.com/goplus/llpkgstore.git
cd llpkgstore

# 构建
go build -o llpkgstore ./cmd/llpkgstore

# 安装到系统路径
sudo mv llpkgstore /usr/local/bin/
```

#### 方法 2: 使用 go install

```bash
# 安装最新版本
go install github.com/PengPengPeng717/llpkgstore@latest
# go install github.com/goplus/llpkgstore/cmd/llpkgstore@latest

# 安装特定版本
go install github.com/PengPengPeng717/llpkgstore@v1.0.0
# go install github.com/goplus/llpkgstore/cmd/llpkgstore@v1.0.0
```


### 环境配置

#### 环境变量

```bash
# LLGo 安装路径
export LLGO_ROOT=/path/to/llgo

# Python 环境路径
export PYTHONHOME=/path/to/python

# GitHub API 访问令牌
export GITHUB_TOKEN=your_github_token

# Go 模块代理 (可选)
export GOPROXY=https://goproxy.cn,direct
```

#### 配置文件

创建 `~/.llpkgstore/config.yaml`:

```yaml
# 默认配置
defaults:
  type: "cpp"  # 默认包类型
  installer: "conan"  # 默认安装器

# 工具路径
tools:
  llcppg: "/usr/local/bin/llcppg"
  llpyg: "/usr/local/bin/llpyg"
  conan: "/usr/local/bin/conan"
  pip: "/usr/local/bin/pip3"

# 输出目录
output:
  base_dir: "./output"
  demo_dir: "./_demo"
```

## 🚀 快速开始

### 创建 C/C++ 包

#### 1. 创建目录结构

```bash
mkdir mylib
cd mylib
```

#### 2. 创建配置文件

创建 `llpkg.cfg`:

```json
{
  "upstream": {
    "installer": {
      "name": "conan",
      "config": {
        "options": ""
      }
    },
    "package": {
      "name": "mylib",
      "version": "1.0.0"
    }
  },
  "name": "mylib",
  "include": ["mylib.h"],
  "cflags": "$(pkg-config --cflags mylib)",
  "libs": "$(pkg-config --libs mylib)"
}
```

#### 3. 生成绑定

```bash
llpkgstore generate
```

#### 4. 测试包

```bash
llpkgstore demotest
```

### 创建 Python 包

#### 1. 创建目录结构

```bash
mkdir mypythonlib
cd mypythonlib
```

#### 2. 创建配置文件

创建 `llpkg.cfg`:

```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip",
      "config": {
        "options": ""
      }
    },
    "package": {
      "name": "requests",
      "version": "2.31.0"
    }
  }
}
```

#### 3. 安装依赖

```bash
llpkgstore install llpkg.cfg
```

#### 4. 生成绑定

```bash
llpkgstore generate
```

#### 5. 测试包

```bash
llpkgstore demotest
```

## 📚 详细使用指南

### 包类型检测

llpkgstore 会自动检测当前目录的包类型：

- **C/C++ 包**: `llpkg.cfg` 中 `type` 字段为空或 `"cpp"`
- **Python 包**: `llpkg.cfg` 中 `type` 字段为 `"python"`

### 工作流程

#### 1. 配置阶段
- 创建 `llpkg.cfg` 配置文件
- 定义包的基本信息和依赖关系
- 配置安装器和包管理器选项

#### 2. 安装阶段
- 使用指定的安装器下载上游包
- 安装依赖和二进制文件
- 准备生成环境

#### 3. 生成阶段
- 调用相应的生成工具 (`llcppg` 或 `llpyg`)
- 生成 Go 绑定代码
- 创建 Go 模块文件

#### 4. 测试阶段
- 生成演示代码
- 运行测试验证功能
- 检查生成的包质量

#### 5. 发布阶段
- 创建 Git 标签
- 上传到 GitHub Releases
- 更新版本记录

## ⚙️ 配置说明

### 配置文件结构

#### 通用字段

| 字段 | 类型 | 默认值 | 必需 | 描述 |
|------|------|--------|------|------|
| `type` | `string` | `"cpp"` | ❌ | 包类型: `"python"` 或 `"cpp"` |
| `name` | `string` | - | ❌ | 包名称 (可选，从目录名推断) |

#### 上游配置

```json
{
  "upstream": {
    "installer": {
      "name": "conan|pip",
      "config": {
        "options": "安装器特定选项"
      }
    },
    "package": {
      "name": "包名",
      "version": "版本号"
    }
  }
}
```

#### C/C++ 特定配置

```json
{
  "cflags": "编译标志",
  "libs": "链接标志",
  "include": ["头文件路径"],
  "trimPrefixes": ["前缀列表"],
  "cplusplus": false
}
```

#### Python 特定配置

```json
{
  "type": "python",
  "llpyg": {
    "module_name": "自定义模块名",
    "extract_depth": 3,
    "skip_private": true
  }
}
```

### 配置验证

llpkgstore 会自动验证配置文件：

- 检查必需字段
- 验证字段类型和格式
- 检查包名称一致性
- 验证安装器配置

## 📋 命令参考

### 核心命令

#### `llpkgstore generate`

生成 Go 语言绑定。

**用法**:
```bash
llpkgstore generate [flags]
```

**功能**:
- 检测包类型
- 安装上游依赖
- 生成 Go 绑定代码
- 创建演示代码

**标志**:
- `--force`: 强制重新生成
- `--verbose`: 详细输出
- `--output`: 指定输出目录

#### `llpkgstore install`

安装上游依赖包。

**用法**:
```bash
llpkgstore install <config_file> [flags]
```

**功能**:
- 解析配置文件
- 下载上游包
- 安装依赖
- 准备生成环境

**标志**:
- `-o, --output`: 指定输出目录
- `--force`: 强制重新安装

#### `llpkgstore verification`

验证 PR 和包配置。

**用法**:
```bash
llpkgstore verification [flags]
```

**功能**:
- 验证配置文件正确性
- 检查目录名称一致性
- 验证 PR 提交信息格式
- 运行基本测试

**环境变量**:
- `GITHUB_TOKEN`: GitHub API 访问令牌

#### `llpkgstore release`

构建和上传二进制包。

**用法**:
```bash
llpkgstore release [flags]
```

**功能**:
- 构建多平台二进制包
- 创建 Git 标签
- 上传到 GitHub Releases
- 更新版本记录

**环境变量**:
- `GITHUB_TOKEN`: GitHub API 访问令牌

#### `llpkgstore postprocessing`

处理合并的 PR。

**用法**:
```bash
llpkgstore postprocessing [flags]
```

**功能**:
- 创建 Git 标签
- 更新版本记录
- 清理遗留分支
- 创建 GitHub Release

**环境变量**:
- `GITHUB_TOKEN`: GitHub API 访问令牌

#### `llpkgstore demotest`

运行演示程序。

**用法**:
```bash
llpkgstore demotest [flags]
```

**功能**:
- 运行 `_demo/` 目录下的演示程序
- 验证生成的 Go 绑定功能
- 生成测试报告

**标志**:
- `--timeout`: 测试超时时间
- `--verbose`: 详细输出

### 辅助命令

#### `llpkgstore test`

测试 Python 验证功能。

**用法**:
```bash
llpkgstore test [flags]
```

**功能**:
- 运行 Python 包验证测试
- 检查环境配置
- 验证工具可用性

#### `llpkgstore labelcreate`

创建版本维护标签。

**用法**:
```bash
llpkgstore labelcreate -l <label> [flags]
```

**功能**:
- 创建版本分支策略
- 管理维护标签
- 支持 C/C++ 和 Python 包

#### `llpkgstore issueclose`

清理未使用的分支。

**用法**:
```bash
llpkgstore issueclose [flags]
```

**功能**:
- 当 Issue 关闭时清理相关资源
- 删除不再需要的分支
- 维护仓库整洁

## ❓ 常见问题

### 安装问题

#### Q: 安装后找不到 `llpkgstore` 命令

**A**: 检查以下几点：
1. 确保 `$GOPATH/bin` 在 `PATH` 环境变量中
2. 重新安装: `go install github.com/goplus/llpkgstore/cmd/llpkgstore@latest`
3. 检查安装路径: `which llpkgstore`

#### Q: 依赖工具缺失错误

**A**: 安装必需的依赖工具：
```bash
# 安装 LLGo
go install github.com/goplus/llgo@latest

# 安装 LLPyg
go install github.com/toaction/llpyg/cmd/llpyg@latest

# 安装 Conan (C++ 包)
pip install conan
```

### 配置问题

#### Q: 配置文件解析错误

**A**: 检查 JSON 格式：
1. 验证 JSON 语法正确性
2. 检查必需字段是否存在
3. 确保字段类型正确

#### Q: 包类型检测失败

**A**: 检查配置：
1. 确保 `llpkg.cfg` 文件存在
2. 检查 `type` 字段值是否正确
3. 验证目录名称与包名称一致

### 生成问题

#### Q: 生成过程中断

**A**: 常见原因和解决方案：
1. **网络问题**: 检查网络连接，配置代理
2. **权限问题**: 确保有写入权限
3. **依赖缺失**: 安装缺失的依赖包
4. **工具错误**: 检查生成工具状态

#### Q: 生成的代码无法编译

**A**: 检查生成过程：
1. 运行 `llpkgstore demotest` 验证
2. 检查 Go 版本兼容性
3. 查看错误日志定位问题

### 发布问题

#### Q: GitHub Release 创建失败

**A**: 检查环境：
1. 确保 `GITHUB_TOKEN` 环境变量设置正确
2. 检查仓库权限
3. 验证标签格式

#### Q: 版本冲突

**A**: 解决版本冲突：
1. 检查现有版本标签
2. 使用 `--force` 标志强制更新
3. 手动清理冲突的标签

## 🔧 故障排除

### 调试模式

启用详细日志输出：

```bash
# 设置环境变量
export DEBUG=1
export VERBOSE=1

# 执行命令
llpkgstore generate --verbose
```

### 日志文件

llpkgstore 会生成详细的日志：

- **运行日志**: `~/.llpkgstore/logs/runtime.log`
- **错误日志**: `~/.llpkgstore/logs/error.log`
- **调试日志**: `~/.llpkgstore/logs/debug.log`

### 常见错误码

| 错误码 | 描述 | 解决方案 |
|--------|------|----------|
| `E001` | 配置文件未找到 | 检查当前目录和文件路径 |
| `E002` | 配置文件格式错误 | 验证 JSON 语法和字段 |
| `E003` | 依赖工具缺失 | 安装相应的工具 |
| `E004` | 权限不足 | 检查文件和目录权限 |
| `E005` | 网络连接失败 | 检查网络和代理设置 |
| `E006` | 包类型不支持 | 检查 `type` 字段值 |
| `E007` | 版本冲突 | 解决版本冲突或使用 `--force` |


## 🤝 贡献指南

### 开发环境设置

#### 1. 克隆仓库

```bash
git clone https://github.com/PengPengPeng717/llpkgstore.git
# git clone https://github.com/goplus/llpkgstore.git
cd llpkgstore
```

#### 2. 安装依赖

```bash
go mod download
go mod tidy
```

#### 3. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./cmd/llpkgstore/internal_python

# 运行基准测试
go test -bench=. ./...
```

### 代码规范

#### Go 代码规范

- 遵循 [Go 官方代码规范](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化代码
- 添加适当的注释和文档
- 编写单元测试

#### 提交规范

使用语义化提交信息：

- `feat`: 新功能
- `fix`: 修复问题
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

#### 示例

```bash
git commit -m "feat: add support for Python package type detection"
git commit -m "fix: resolve version mapping conflict in postprocessing"
git commit -m "docs: update installation guide for macOS users"
```

### 测试指南

#### 单元测试

```bash
# 运行特定包的测试
go test ./cmd/llpkgstore/internal_python -v

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### 集成测试

```bash
# 运行集成测试
go test -tags=integration ./...

# 运行端到端测试
go test -tags=e2e ./...
```

#### 性能测试

```bash
# 运行基准测试
go test -bench=. -benchmem ./...

# 运行性能分析
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

## 🔄 版本管理

### 统一版本管理机制

llpkgstore v2.0+ 引入了统一的版本管理机制，为 C/C++ 和 Python 包提供一致的版本处理体验。

#### 版本提取流程

1. **提交消息解析**: 从最新的提交消息中提取版本信息
2. **格式支持**: 支持多种版本格式
   - `Release-as: package_name/vX.X.X`
   - `Release: vX.X.X`
   - `Version: vX.X.X`
3. **Git 标签回退**: 如果提交消息中未找到版本，自动从 Git 标签获取

#### 版本记录机制

**双重版本记录**:
- **本地记录**: 包目录下的 `llpkgstore.json`
- **集中记录**: `llpkg/public/llpkgstore.json`

**版本映射格式**:
```json
{
  "packages": {
    "package_name": {
      "versions": [
        {
          "python": "0.9.0",
          "go": ["v8.0.0", "v9.0.0", "v10.0.0"]
        }
      ]
    }
  }
}
```

#### 自动 Git 标签

- **标签创建**: 基于提取的版本信息自动创建 Git 标签
- **标签推送**: 自动推送到远程仓库
- **重复检测**: 检查标签是否已存在，避免重复创建

#### 版本验证

- **格式验证**: 确保版本号符合语义化版本规范
- **冲突检测**: 检测版本冲突和重复
- **回滚机制**: 支持版本回滚和修复

### 发布流程

#### 1. 版本管理

```bash
# 创建新版本标签
git tag -a v1.1.0 -m "Release version 1.1.0"
git push origin v1.1.0
```

#### 2. 构建发布

```bash
# 构建多平台二进制文件
make build-all

# 创建 GitHub Release
gh release create v1.1.0 --title "Release v1.1.0" --notes "Release notes"
```

#### 3. 文档更新

- 更新 `CHANGELOG.md`
- 更新版本号
- 更新依赖版本

## 📝 更新日志

### [v2.0.0] - 2024-01-XX

#### 🎉 重大更新
- **Python 包支持**: 完整支持 Python 包的 Go 绑定生成
- **智能路由**: 根据包类型自动选择处理逻辑
- **统一配置**: 支持 C/C++ 和 Python 包的统一配置格式

#### ✨ 新功能
- 新增 `llpkgstore generate` 命令的 Python 支持
- 新增 `llpkgstore install` 命令的 pip 集成
- 新增 `llpkgstore demotest` 命令的 Python 包测试
- 新增 `llpkgstore test` 命令用于 Python 验证

#### 🔧 改进
- 重构主程序架构，支持动态命令路由
- 优化错误处理和用户提示
- 改进配置文件验证和错误报告
- 增强 CI/CD 集成支持

#### 🐛 修复
- 修复环境变量冲突问题
- 修复配置文件解析错误
- 修复版本映射逻辑问题
- 修复 GitHub Actions 工作流问题

### [v1.0.0] - 2023-XX-XX

#### 🎉 初始版本
- **C/C++ 包支持**: 完整支持 C/C++ 库的 Go 绑定生成
- **Conan 集成**: 支持 Conan 包管理器
- **版本管理**: 复杂的语义化版本映射
- **CI/CD 集成**: GitHub Actions 工作流支持

#### ✨ 核心功能
- `llpkgstore generate`: 生成 C/C++ 库的 Go 绑定
- `llpkgstore install`: 安装 C/C++ 依赖
- `llpkgstore verification`: 验证包配置
- `llpkgstore release`: 发布二进制包
- `llpkgstore postprocessing`: 处理合并的 PR

### 版本规划

#### 即将发布
- **v2.1.0**: 性能优化和错误处理改进
- **v2.2.0**: 配置模板和验证器增强

#### 长期规划
- **v3.0.0**: Rust 包支持
- **v4.0.0**: Node.js 包支持
- **v5.0.0**: Java 包支持

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🔗 相关链接

- [项目仓库](https://github.com/goplus/llpkgstore)
- [LLGo 项目](https://github.com/goplus/llgo)
- [LLPyg 项目](https://github.com/toaction/llpyg)
- [LLCppg 项目](https://github.com/goplus/llcppg)
- [文档网站](https://goplus.github.io/llpkg/)
- [问题反馈](https://github.com/goplus/llpkgstore/issues)
- [讨论区](https://github.com/goplus/llpkgstore/discussions)

## 🙏 致谢

感谢所有为 llpkgstore 项目做出贡献的开发者和用户。特别感谢：

- [LLGo 团队](https://github.com/goplus/llgo) 提供的基础框架
- [LLPyg 团队](https://github.com/toaction/llpyg) 提供的 Python 绑定支持
- [LLCppg 团队](https://github.com/goplus/llcppg) 提供的 C/C++ 绑定支持
- 所有提交 Issue 和 Pull Request 的贡献者

---

**最后更新**: 2024年1月

如有问题或建议，请通过 [GitHub Issues](https://github.com/goplus/llpkgstore/issues) 联系我们。
