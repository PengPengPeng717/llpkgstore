# llpkgstore - 统一包分发服务

[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/llpkgstore)](https://goreportcard.com/report/github.com/goplus/llpkgstore)
[![Go Version](https://img.shields.io/github/go-mod/go-version/goplus/llpkgstore)](https://github.com/goplus/llpkgstore/blob/main/go.mod)
[![License](https://img.shields.io/github/license/goplus/llpkgstore)](https://github.com/goplus/llpkgstore/blob/main/LICENSE)
[![Build Status](https://github.com/goplus/llpkgstore/workflows/CI/badge.svg)](https://github.com/goplus/llpkgstore/actions)

> **llpkgstore** 是一个专为 [**LLGo**](https://github.com/goplus/llgo) 设计的综合包分发服务，为多语言生态系统提供可信赖且便捷的语言绑定访问。

## 📋 目录

- [概述](#概述)
- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [使用指南](#使用指南)
- [配置参考](#配置参考)
- [API 参考](#api-参考)
- [架构设计](#架构设计)
- [开发指南](#开发指南)
- [常见问题](#常见问题)
- [更新日志](#更新日志)
- [许可证](#许可证)

## 概述

llpkgstore 是一个专为 [**LLGo**](https://github.com/goplus/llgo) 设计的统一包分发服务，为多语言生态系统提供可信赖且便捷的语言绑定访问。

### 什么是 llpkgstore？

llpkgstore 是一个综合性的包管理工具，它能够：

- **自动生成语言绑定**: 为 C/C++ 和 Python 库自动生成 Go 语言绑定
- **统一版本管理**: 提供一致的版本映射和管理机制
- **智能包检测**: 优先使用系统环境中的包，提高效率
- **CI/CD 集成**: 与 GitHub Actions 无缝集成，自动化包生成和发布

### 核心价值

- **🚀 简化集成**: 将复杂的跨语言绑定生成过程简化为几个命令
- **🔒 安全可靠**: 通过自动化流程和验证机制确保包质量
- **🌍 多语言支持**: 统一支持 C/C++ 和 Python 生态系统
- **⚡ 高效便捷**: 智能检测和缓存机制，减少重复工作

## 特性

### 多语言支持

llpkgstore 目前支持以下编程语言和平台：

| 语言/平台 | 状态 | 生成工具 | 包管理器 | 主要特性 |
|-----------|------|----------|----------|----------|
| **C/C++** | ✅ 完全支持 | `llcppg` | Conan | 二进制分发、头文件、.pc 文件 |
| **Python** | ✅ 完全支持 | `llpyg` | pip | 模块绑定、Go 接口、类型安全 |
| **JavaScript** | 🚧 计划中 | - | - | 计划支持 |
| **Rust** | 🚧 计划中 | - | - | 计划支持 |

### 核心功能

#### 🔄 统一版本管理
- **智能版本提取**: 从提交消息自动提取版本信息
- **版本格式验证**: 支持语义化版本和自定义格式
- **双重记录机制**: 本地和集中式版本记录同步
- **自动 Git 标签**: 基于版本信息自动创建和推送标签

#### 🎯 智能包检测
- **系统包优先**: 优先使用系统环境中已安装的包
- **临时安装**: 仅在必要时进行临时包安装
- **缓存机制**: 智能缓存减少重复下载和安装

#### 🚀 自动化工作流
- **CI/CD 集成**: 与 GitHub Actions 无缝集成
- **自动发布**: 自动创建 GitHub Releases
- **测试验证**: 自动生成和运行测试代码
- **错误处理**: 完善的错误处理和回滚机制

#### 🔧 开发者体验
- **统一 CLI**: 一致的命令行接口
- **配置驱动**: 基于配置文件的灵活设置
- **详细日志**: 完整的操作日志和调试信息
- **文档生成**: 自动生成使用文档和示例

## 🏗️ 系统架构

### 整体架构图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户界面层     │    │   业务逻辑层     │    │   数据存储层     │
│                │    │                │    │                │
│ • CLI 命令      │◄──►│ • 包管理逻辑     │◄──►│ • 配置文件      │
│ • 配置接口      │    │ • 版本管理      │    │ • 版本记录      │
│ • 错误处理      │    │ • 生成器调用     │    │ • 元数据存储     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 核心组件

#### 1. 命令层
- 统一的 CLI 接口
- 自动包类型检测
- 智能路由和分发

#### 2. 业务逻辑层
- 包管理核心逻辑
- 版本管理和映射
- 生成器集成

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

## 安装

### 系统要求

在安装 llpkgstore 之前，请确保您的系统满足以下要求：

| 组件 | 最低版本 | 说明 |
|------|----------|------|
| **Go** | 1.19+ | 用于构建和运行 llpkgstore |
| **Python** | 3.7+ | 用于 Python 包支持 |
| **Git** | 2.0+ | 用于版本控制和标签管理 |
| **操作系统** | - | Linux、macOS、Windows |

### 安装方法

#### 方法 1: 使用 go install (推荐)

```bash
go install github.com/goplus/llpkgstore/cmd/llpkgstore@latest
```

#### 方法 2: 从源码构建

```bash
# 克隆仓库
git clone https://github.com/goplus/llpkgstore.git
cd llpkgstore

# 构建
go build -o llpkgstore ./cmd/llpkgstore

# 安装到系统路径
sudo cp llpkgstore /usr/local/bin/
```

#### 方法 3: 下载预编译二进制文件

从 [GitHub Releases](https://github.com/goplus/llpkgstore/releases) 下载适合您系统的预编译版本。

### 验证安装

安装完成后，验证 llpkgstore 是否正确安装：

```bash
llpkgstore --version
```

如果安装成功，您应该看到类似以下的输出：

```
llpkgstore version 2.0.0
```

### 环境配置

#### Python 环境配置

llpkgstore 需要 Python 环境来支持 Python 包的处理：

```bash
# 检查 Python 版本
python3 --version

# 检查 pip 版本
pip3 --version

# 可选：创建专用的虚拟环境
python3 -m venv ~/.llpkgstore-env
source ~/.llpkgstore-env/bin/activate  # Linux/macOS
# 或
~/.llpkgstore-env\Scripts\activate     # Windows
```

#### C/C++ 环境配置 (可选)

如果您需要处理 C/C++ 包，可以安装 Conan：

```bash
# 安装 Conan
pip3 install conan

# 配置 Conan
conan profile detect --force
```

### 故障排除

#### 常见安装问题

**问题 1**: `command not found: llpkgstore`
```bash
# 解决方案：确保 Go bin 目录在 PATH 中
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

**问题 2**: Python 版本不兼容
```bash
# 解决方案：使用 pyenv 管理 Python 版本
curl https://pyenv.run | bash
pyenv install 3.9.0
pyenv global 3.9.0
```

## 快速开始

本指南将帮助您快速上手 llpkgstore，通过几个简单的步骤创建一个 Python 包的 Go 绑定。

### 步骤 1: 创建项目目录

```bash
mkdir my-python-package
cd my-python-package
```

### 步骤 2: 创建配置文件

创建 `llpkg.cfg` 配置文件，指定要处理的 Python 包：

```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "requests",
      "version": "2.31.0"
    }
  },
  "llpyg": {
    "output_dir": "./bindings",
    "mod_name": "github.com/your-org/requests",
    "mod_depth": 1
  }
}
```

### 步骤 3: 安装包

```bash
llpkgstore install llpkg.cfg
```

### 步骤 4: 生成 Go 绑定

```bash
llpkgstore generate
```

### 步骤 5: 验证结果

检查生成的文件：

```bash
ls -la bindings/
```

您应该看到生成的 Go 文件，包括：
- `requests.go` - 主要的 Go 绑定文件
- `requests_autogen_link.go` - 自动生成的链接文件

### 步骤 6: 测试绑定

```bash
cd bindings
go mod tidy
go run _demo/main.go
```

## 使用指南

### 基本命令

llpkgstore 提供以下主要命令：

| 命令 | 描述 | 示例 |
|------|------|------|
| `install` | 安装指定的包 | `llpkgstore install llpkg.cfg` |
| `generate` | 生成 Go 绑定 | `llpkgstore generate` |
| `postprocessing` | 后处理（版本管理） | `llpkgstore postprocessing` |
| `release` | 创建发布 | `llpkgstore release` |

### 配置文件格式

#### Python 包配置

```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "package_name",
      "version": "1.0.0"
    }
  },
  "llpyg": {
    "output_dir": "./output",
    "mod_name": "github.com/your-org/package_name",
    "mod_depth": 1
  }
}
```

#### C/C++ 包配置

```json
{
  "type": "cpp",
  "upstream": {
    "installer": {
      "name": "conan"
    },
    "package": {
      "name": "mylib",
      "version": "1.0.0"
    }
  },
  "name": "mylib",
  "version": "1.0.0"
}
```

#### Python 包配置

```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "requests",
      "version": "2.31.0"
    }
  },
  "llpyg": {
    "output_dir": "./test",
    "mod_name": "github.com/PengPengPeng717/llpkg/requests",
    "mod_depth": 1
  }
}
```

### Python 包支持详解

#### 核心特性

- ✅ **完整 Python 包支持**: 为任何 Python 包生成 Go 绑定
- ✅ **类型安全**: 在生成的 Go 代码中保持类型信息
- ✅ **模块深度控制**: 可配置的嵌套模块提取深度
- ✅ **自定义模块名**: 支持自定义 Go 模块名
- ✅ **自动测试**: 生成演示代码进行验证
- ✅ **统一版本管理**: 与 C/C++ 包一致的版本映射
- ✅ **智能版本提取**: 从提交消息自动提取版本信息
- ✅ **双重版本记录**: 本地和集中式版本记录文件
- ✅ **自动 Git 标签**: 基于版本自动创建和推送 Git 标签
- ✅ **CI/CD 集成**: 完整的 GitHub Actions 支持

#### Python 包生成流程

1. **包安装**: 使用 pip 安装指定的 Python 包
2. **绑定生成**: 生成 Go 接口和类型安全的绑定
3. **模块配置**: 配置输出目录、模块名和提取深度
4. **Go 模块创建**: 创建带有依赖关系的正确 Go 模块
5. **测试验证**: 使用演示代码测试生成的绑定

#### llpyg 配置选项

| 字段 | 类型 | 默认值 | 可选 | 描述 |
|------|------|--------|------|------|
| `llpyg.output_dir` | `string` | `"./test"` | ✅ | 生成文件的输出目录 |
| `llpyg.mod_name` | `string` | 包名 | ✅ | Go 模块名 |
| `llpyg.mod_depth` | `int` | `1` | ✅ | 最大模块提取深度 (0-10) |

#### Python 包测试

Python 包在 `_demo` 目录中包含演示代码，用于验证：
- 包导入和编译
- 基本功能测试
- 类型安全验证
- 与 Go 生态系统的集成

### 3. 生成 Go 绑定

```bash
# 安装上游包
llpkgstore install

# 生成 Go 绑定
llpkgstore generate

# 运行测试
llpkgstore test
```

## 📖 详细使用指南

### 命令概览

| 命令 | 描述 | 示例 |
|------|------|------|
| `install` | 安装上游包 | `llpkgstore install` |
| `generate` | 生成 Go 绑定 | `llpkgstore generate` |
| `test` | 运行测试 | `llpkgstore test` |
| `postprocessing` | 后处理 | `llpkgstore postprocessing` |
| `release` | 创建发布 | `llpkgstore release` |
| `verification` | 验证包 | `llpkgstore verification` |

### 配置文件详解

#### 基本结构

```json
{
  "type": "包类型",
  "upstream": {
    "installer": {
      "name": "安装器名称"
    },
    "package": {
      "name": "包名",
      "version": "版本号"
    }
  }
}
```

#### 包类型特定配置

**Python 包**:
```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "包名",
      "version": "版本号"
    }
  },
  "llpyg": {
    "output_dir": "./output",
    "mod_name": "github.com/user/repo",
    "mod_depth": 1
  }
}
```

## 🔄 版本管理

### 统一版本管理机制

llpkgstore v2.0+ 实现了完全统一的版本管理机制，C/C++ 和 Python 包现在使用相同的版本处理逻辑和接口。

#### 统一架构

- **共享接口**: 所有包类型都使用 `DefaultClient` 接口
- **统一版本提取**: 使用相同的版本提取逻辑和格式支持
- **一致的后处理**: 相同的 GitHub Release 创建和 Git 标签管理
- **兼容的版本记录**: 统一的版本记录格式和更新机制

#### 版本提取流程

1. **提交消息解析**: 从最新的提交消息中提取版本信息
2. **格式支持**: 支持多种版本格式
   - `Release-as: package_name/vX.X.X`

3. **Git 标签回退**: 如果提交消息中未找到版本，自动从 Git 标签获取
4. **统一验证**: 使用相同的版本格式验证和冲突检测

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

### 版本映射逻辑详解

#### 核心设计原则

1. **统一版本管理**: Python 和 C++ 包使用相同的版本管理逻辑
2. **基于 Commit 消息的版本提取**: 从 Git commit 消息中提取版本信息
3. **集中式版本记录**: 使用 `llpkgstore.json` 文件记录版本映射

#### 版本提取优先级

1. **最新提交消息** (用于 CI 环境)
2. **Git 标签** (回退机制)
3. **错误处理** (清晰的错误消息)

#### 版本格式验证

支持的版本格式：
- `Release-as: {包名}/v{版本号}` (推荐)
- `Release-as: {包名}/{版本号}`
- 语义化版本号 (SemVer): `v1.2.3`, `1.2.3`

#### 版本记录更新策略

```go
// 更新版本记录文件
func updateLLPkgStoreJSON(packageName, upstreamVersion, mappedVersion string) error {
    // 1. 读取现有版本记录
    data, err := os.ReadFile("llpkgstore.json")
    if err != nil {
        // 文件不存在时创建新的
        data = []byte("{}")
    }
    
    // 2. 解析 JSON
    var versionData LLPkgStoreJSON
    if len(strings.TrimSpace(string(data))) == 0 {
        // 空文件时初始化
        versionData = LLPkgStoreJSON{Packages: make(map[string]PackageInfo)}
    } else {
        err = json.Unmarshal(data, &versionData)
        if err != nil {
            return err
        }
    }
    
    // 3. 更新版本信息
    versionData.Write(packageName, upstreamVersion, mappedVersion)
    
    // 4. 写回文件
    updatedData, err := json.MarshalIndent(versionData, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile("llpkgstore.json", updatedData, 0644)
}
```

## 开发指南

### 开发环境设置

#### 1. 克隆仓库

```bash
git clone https://github.com/goplus/llpkgstore.git
cd llpkgstore
```

#### 2. 安装依赖

```bash
# 安装 Go 依赖
go mod tidy

# 安装开发工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/goplus/llgo@latest
```

#### 3. 构建项目

```bash
# 构建 llpkgstore
go build -o llpkgstore ./cmd/llpkgstore

# 运行测试
go test ./...

# 运行 lint 检查
golangci-lint run
```

### 代码结构

#### 核心模块

- **`cmd/llpkgstore/`**: 命令行接口实现
- **`internal/actions/`**: 核心业务逻辑
- **`config/`**: 配置管理和验证
- **`upstream/`**: 上游包管理器集成
- **`metadata/`**: 元数据管理和缓存

#### 开发规范

1. **代码风格**: 遵循 Go 官方代码规范
2. **测试覆盖**: 新功能必须包含单元测试
3. **文档更新**: 修改功能时同步更新文档
4. **错误处理**: 提供清晰的错误信息和处理机制

### 贡献流程

#### 1. 创建 Issue

在提交代码之前，请先创建 Issue 描述您要解决的问题或添加的功能。

#### 2. Fork 和分支

```bash
# Fork 仓库后克隆
git clone https://github.com/your-username/llpkgstore.git
cd llpkgstore

# 创建功能分支
git checkout -b feature/your-feature-name
```

#### 3. 开发和测试

```bash
# 开发功能
# ... 编写代码 ...

# 运行测试
go test ./...

# 运行 lint
golangci-lint run

# 构建验证
go build ./cmd/llpkgstore
```

#### 4. 提交代码

```bash
# 添加更改
git add .

# 提交更改
git commit -m "feat: add your feature description"

# 推送分支
git push origin feature/your-feature-name
```

#### 5. 创建 Pull Request

在 GitHub 上创建 Pull Request，包含：

- 清晰的标题和描述
- 相关的 Issue 链接
- 测试结果截图（如适用）
- 文档更新说明

### 发布流程

#### 版本管理

llpkgstore 使用语义化版本控制：

- **主版本号**: 不兼容的 API 修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

#### 发布步骤

1. **更新版本号**: 修改 `go.mod` 中的版本
2. **更新文档**: 更新 CHANGELOG.md
3. **创建标签**: `git tag v1.0.0`
4. **推送标签**: `git push origin v1.0.0`
5. **创建 Release**: 在 GitHub 上创建 Release

### 调试指南

#### 启用调试模式

```bash
# 设置调试环境变量
export LLPKGSTORE_DEBUG=1

# 运行命令查看详细日志
llpkgstore generate --verbose
```

#### 常见调试场景

1. **包安装失败**: 检查网络连接和包管理器配置
2. **生成器错误**: 验证配置文件格式和依赖项
3. **版本冲突**: 检查版本映射和 Git 标签

## ❓ 常见问题

### 安装问题

**Q: 安装后无法找到 llpkgstore 命令**
A: 确保 `/usr/local/bin` 在您的 PATH 环境变量中，或使用 `go install` 安装到 GOPATH。

**Q: Python 包安装失败**
A: 确保已安装 Python 3.8+ 和 pip，并检查网络连接。

### 配置问题

**Q: 配置文件格式错误**
A: 使用 JSON 验证器检查配置文件格式，确保所有必需字段都存在。

**Q: 包类型检测失败**
A: 确保 `type` 字段设置为 `"python"` 或 `"cpp"`。

### 生成问题

**Q: Go 绑定生成失败**
A: 检查上游包是否正确安装，确保 llpyg 或 llcppg 工具可用。

## 🔧 故障排除

### 调试模式

```bash
# 启用详细输出
llpkgstore --verbose generate

# 检查配置
llpkgstore config validate
```

### 日志文件

- **安装日志**: `~/.llpkgstore/logs/install.log`
- **生成日志**: `~/.llpkgstore/logs/generate.log`
- **错误日志**: `~/.llpkgstore/logs/error.log`

### 常见错误

#### 1. 包未找到
```
Error: package not found
```
**解决方案**: 检查包名和版本号是否正确。

#### 2. 配置错误
```
Error: invalid configuration
```
**解决方案**: 验证配置文件格式和必需字段。

#### 3. 权限问题
```
Error: permission denied
```
**解决方案**: 检查文件权限和目录访问权限。

## 🤝 贡献指南

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/goplus/llpkgstore.git
cd llpkgstore

# 安装依赖
go mod download

# 运行测试
go test ./...
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

# 运行性能分析
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

## 📝 更新日志

### v2.0.0 (最新)

#### 🎉 重大更新
- **统一版本管理**: C/C++ 和 Python 包使用一致的版本管理机制
- **智能版本提取**: 从提交消息自动提取版本信息
- **双重版本记录**: 同时维护本地和集中式版本记录文件
- **自动 Git 标签**: 基于版本信息自动创建和推送 Git 标签

#### ✨ 新功能
- **Python 包支持**: 完整的 Python 包 Go 绑定生成
- **统一命令接口**: 自动包类型检测和智能路由
- **增强错误处理**: 结构化的错误处理和恢复机制
- **CI/CD 集成**: 优化的 GitHub Actions 工作流

#### 🔧 改进
- **架构统一**: 消除了 C++ 和 Python 部分的架构不一致
- **代码复用**: 减少了重复代码，提高了维护性
- **用户体验**: 一致的命令接口和错误消息格式

### v1.0.0

#### 🎉 初始版本
- **C/C++ 包支持**: 完整支持 C/C++ 库的 Go 绑定生成
- **Conan 集成**: 支持 Conan 包管理器
- **版本管理**: 复杂的语义化版本映射
- **CI/CD 集成**: GitHub Actions 工作流支持

#### ✨ 核心功能
- **包管理**: 自动包安装和依赖管理
- **绑定生成**: 使用 llcppg 生成 Go 绑定
- **测试框架**: 自动生成演示代码和测试
- **发布管理**: GitHub Releases 集成

---

## 📚 相关资源

- **[架构文档](./ARCHITECTURE.md)**: 详细的系统架构说明
- **[技术文档](./llpkgstore.md)**: 详细的技术设计和实现细节
- **[GitHub 仓库](https://github.com/goplus/llpkgstore)**: 源代码和问题跟踪
- **[LLGo 项目](https://github.com/goplus/llgo)**: LLGo 语言扩展

---

**llpkgstore** - 让跨语言开发更简单、更可靠！ 🚀
