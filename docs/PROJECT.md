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
- [版本管理](#版本管理)
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

## 🚀 安装与部署

### 系统要求

- **Go**: 1.21 或更高版本
- **操作系统**: Linux, macOS, Windows
- **依赖工具**: 
  - Python 3.8+ (用于 Python 包)
  - Conan (用于 C/C++ 包)
  - Git (用于版本管理)

### 安装方法

#### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/goplus/llpkgstore.git
cd llpkgstore

# 构建
go build -o llpkgstore ./cmd/llpkgstore

# 安装到系统路径
sudo mv llpkgstore /usr/local/bin/
```

#### 使用 Go 安装

```bash
go install github.com/goplus/llpkgstore@latest
```

### 验证安装

```bash
llpkgstore --version
```

## 🚀 快速开始

### 1. 创建包目录

```bash
mkdir my-package
cd my-package
```

### 2. 创建配置文件

创建 `llpkg.cfg` 文件：

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
  }
}
```

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
- **[Python 支持](./python-support.md)**: Python 包支持详细指南
- **[llpyg 配置](./llpyg-config.md)**: llpyg 工具配置说明
- **[GitHub 仓库](https://github.com/goplus/llpkgstore)**: 源代码和问题跟踪
- **[LLGo 项目](https://github.com/goplus/llgo)**: LLGo 语言扩展

---

**llpkgstore** - 让跨语言开发更简单、更可靠！ 🚀
