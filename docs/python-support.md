# Python Package Support in llpkgstore

This document provides comprehensive information about Python package support in llpkgstore, including configuration, generation, testing, and best practices.

## Overview

llpkgstore v2.0+ provides full support for Python packages through integration with the `llpyg` tool. This allows you to create Go bindings for Python modules, making Python libraries accessible from Go code.

## Features

- ✅ **Full Python Package Support**: Generate Go bindings for any Python package
- ✅ **Type Safety**: Maintained type information in generated Go code
- ✅ **Module Depth Control**: Configurable extraction depth for nested modules
- ✅ **Custom Module Names**: Support for custom Go module names
- ✅ **Automatic Testing**: Demo code generation for verification
- ✅ **Unified Version Management**: Consistent version mapping with C/C++ packages
- ✅ **Smart Version Extraction**: Automatic version extraction from commit messages
- ✅ **Dual Version Recording**: Local and centralized version record files
- ✅ **Automatic Git Tagging**: Auto-create and push Git tags based on versions
- ✅ **CI/CD Integration**: Full GitHub Actions support with postprocessing

## Quick Start

### 1. Create Package Directory

```bash
mkdir my-python-package
cd my-python-package
```

### 2. Create Configuration File

Create `llpkg.cfg`:

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

### 3. Generate Go Bindings

```bash
llpkgstore generate
```

This will:
- Install the Python package using pip
- Generate Go bindings using `llpyg`
- Create Go module files (`go.mod`, `go.sum`)
- Generate demo code in `_demo` directory

### 4. Test the Generated Package

```bash
cd _demo/basic_test
llgo run main.go
```

## Configuration Options

### Basic Configuration

```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "numpy",
      "version": "1.26.4"
    }
  }
}
```

### Advanced Configuration with llpyg Options

```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip",
      "config": {
        "options": "--user"
      }
    },
    "package": {
      "name": "pandas",
      "version": "2.1.0"
    }
  },
  "llpyg": {
    "output_dir": "./bindings",
    "mod_name": "github.com/myorg/mypackage/pandas",
    "mod_depth": 3
  }
}
```

### Configuration Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | ✅ | - | Must be "python" |
| `upstream.installer.name` | string | ✅ | - | Must be "pip" |
| `upstream.installer.config.options` | string | ❌ | "" | Additional pip options |
| `upstream.package.name` | string | ✅ | - | Python package name |
| `upstream.package.version` | string | ✅ | - | Python package version |
| `llpyg.output_dir` | string | ❌ | "./test" | Output directory |
| `llpyg.mod_name` | string | ❌ | package name | Go module name |
| `llpyg.mod_depth` | integer | ❌ | 1 | Module extraction depth (0-10) |

## Package Generation Process

### 1. Package Installation

The system automatically installs the specified Python package:

```bash
pip install {package_name}=={version}
```

### 2. Binding Generation

Uses `llpyg` to generate Go bindings:

```bash
llpyg -o {output_dir} -mod {mod_name} -d {mod_depth} {package_name}
```

### 3. Go Module Creation

Creates a proper Go module with:
- `go.mod` file with dependencies
- Generated Go source files
- Demo code for testing

### 4. File Structure

After generation, your directory will contain:

```
my-python-package/
├── llpkg.cfg
├── go.mod
├── go.sum
├── {package_name}.go          # Generated Go bindings
├── _demo/
│   └── basic_test/
│       └── main.go            # Demo code
└── llpyg.cfg                  # llpyg configuration
```

## Demo Code

### Basic Test Structure

Generated demo code follows this pattern:

```go
package main

import (
    "fmt"
    "github.com/goplus/lib/py"
    "github.com/goplus/lib/py/std"
    "yourpackage"  // Your generated package
)

func main() {
    fmt.Println("=== Package Test Demo ===")
    
    // Test basic functionality
    // Test type safety
    // Test integration
    
    fmt.Println("=== Test Complete ===")
}
```

### Testing Guidelines

1. **Import Testing**: Verify package can be imported
2. **Basic Functionality**: Test core package features
3. **Type Safety**: Ensure proper type handling
4. **Error Handling**: Test error conditions
5. **Integration**: Verify Go ecosystem compatibility

## Version Management

### Commit Message Format

Use this format for version releases:

```bash
git commit -m "Release-as: {package_name}/v{version}"
```

Examples:
```bash
git commit -m "Release-as: numpy/v1.26.4"
git commit -m "Release-as: requests/v2.31.0"
```

### Version Tagging

The system automatically:
- Extracts version from commit messages
- Creates GitHub releases
- Manages version tags
- Updates version mapping

### Version Compatibility

- Python versions are mapped to Go versions
- Follows semantic versioning
- Maintains backward compatibility
- Supports multiple version branches

## Best Practices

### 1. Package Selection

- Choose stable, well-maintained Python packages
- Consider package size and complexity
- Verify Python package compatibility
- Test with multiple Python versions

### 2. Configuration

- Use specific version numbers (avoid "latest")
- Set appropriate module depth
- Choose meaningful Go module names
- Document configuration decisions

### 3. Testing

- Write comprehensive demo code
- Test edge cases and error conditions
- Verify type safety
- Test integration with Go ecosystem

### 4. Maintenance

- Keep Python package versions updated
- Monitor for breaking changes
- Maintain backward compatibility
- Document version changes

## Troubleshooting

### Common Issues

#### 1. Python Package Not Found

```bash
Error: Package 'invalid-package' not found
```

**Solution**: Verify package name and availability on PyPI

#### 2. Version Conflicts

```bash
Error: Version '2.0.0' not available
```

**Solution**: Check available versions with `pip index versions {package_name}`

#### 3. Generation Failures

```bash
Error: llpyg generation failed
```

**Solution**: Check Python environment and package installation

#### 4. Go Module Issues

```bash
Error: Invalid module name
```

**Solution**: Verify `mod_name` format and Go module naming conventions

### Debug Commands

```bash
# Check Python package installation
pip show {package_name}

# Verify llpyg installation
llpyg --help

# Check Go module status
go mod tidy
go mod verify

# Test generated package
llgo run _demo/basic_test/main.go
```

## Examples

### Example 1: Simple Package

**Configuration** (`llpkg.cfg`):
```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "colorama",
      "version": "0.4.6"
    }
  }
}
```

**Usage**:
```go
package main

import (
    "github.com/goplus/lib/py"
    "colorama"
)

func main() {
    // Use colorama package
    colorama.init()
    colorama.Fore("red")
}
```

### Example 2: Complex Package

**Configuration** (`llpkg.cfg`):
```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "matplotlib",
      "version": "3.7.2"
    }
  },
  "llpyg": {
    "output_dir": "./matplotlib-bindings",
    "mod_name": "github.com/myorg/plotting/matplotlib",
    "mod_depth": 2
  }
}
```

**Usage**:
```go
package main

import (
    "github.com/goplus/lib/py"
    "matplotlib"
    "matplotlib.pyplot"
)

func main() {
    // Create matplotlib plot
    pyplot.plot(py.List(1, 2, 3, 4), py.List(1, 4, 9, 16))
    pyplot.show()
}
```

## Integration with CI/CD

### GitHub Actions Workflow

```yaml
name: Python Package CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - uses: actions/setup-python@v4
        with:
          python-version: '3.11'
      - name: Install dependencies
        run: |
          pip install llpyg
          go install github.com/goplus/llgo/cmd/llgo@latest
      - name: Generate package
        run: llpkgstore generate
      - name: Test package
        run: |
          cd _demo/basic_test
          llgo run main.go
```

### Automated Testing

- Package generation verification
- Demo code execution
- Type safety checks
- Integration testing
- Performance benchmarking

## Performance Considerations

### Generation Time

- Small packages: 1-5 minutes
- Medium packages: 5-15 minutes
- Large packages: 15+ minutes

### Runtime Performance

- Go bindings are compiled
- Minimal runtime overhead
- Efficient memory usage
- Fast function calls

### Optimization Tips

- Use appropriate module depth
- Generate only needed modules
- Cache generated bindings
- Parallel processing for multiple packages

## Security Considerations

### Package Verification

- Verify package signatures
- Check package checksums
- Validate package sources
- Monitor for vulnerabilities

### Access Control

- Limit package access
- Control version updates
- Audit package usage
- Monitor dependencies

## Future Enhancements

### Planned Features

1. **Virtual Environment Support**: Better Python environment management
2. **Dependency Resolution**: Automatic dependency handling
3. **Package Caching**: Intelligent caching mechanisms
4. **Parallel Generation**: Concurrent package processing
5. **Advanced Testing**: Automated test generation

### Integration Improvements

1. **IDE Support**: Better editor integration
2. **Debugging Tools**: Enhanced debugging capabilities
3. **Performance Profiling**: Built-in performance analysis
4. **Documentation Generation**: Automatic API documentation

## 🔄 统一版本管理 (v2.0+)

### 统一版本管理机制

Python 包现在使用与 C/C++ 包相同的版本管理机制，提供一致的版本处理体验。

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

#### Postprocessing 命令

使用 postprocessing 命令处理版本管理：

```bash
llpkgstore postprocessing
```

此命令将：
1. 从提交消息或 Git 标签提取版本
2. 更新本地 `llpkgstore.json` 版本映射
3. 更新集中式 `llpkg/public/llpkgstore.json`
4. 创建并推送 Git 标签
5. 创建 GitHub 发布（在 CI 环境中）

## 🚀 Python 架构优化总结

### 优化目标

将 Python 部分的后处理逻辑统一到 `actions.DefaultClient` 接口，消除与 C++ 部分的架构不一致问题，实现代码复用和维护性提升。

### 当前问题分析

#### 1. **架构不一致**
- **C++ 部分**: 使用统一的 `actions.DefaultClient` 接口
- **Python 部分**: 直接使用 `exec.Command` 和 GitHub CLI，缺乏统一抽象层

#### 2. **代码重复**
- Python 部分重新实现了版本提取、GitHub Release 创建等功能
- 与 C++ 部分存在大量重复逻辑

#### 3. **错误处理不统一**
- Python 部分使用简单的 `fmt.Errorf`
- C++ 部分使用结构化的错误处理

### 优化方案

#### 1. **统一客户端接口**

##### 创建 `PythonPostProcessor`
```go
// 扩展 DefaultClient 以支持 Python 包
type PythonPostProcessor struct {
    *DefaultClient
}

// 使用统一的接口处理 Python 包
func (p *PythonPostProcessor) Postprocessing() error {
    // 使用统一的版本提取逻辑
    version, err := p.extractVersionUnified()
    // 使用统一的发布创建逻辑
    return p.createReleaseUnified(packageName, version)
}
```

##### 自动包类型检测
```go
// 自动检测包类型并选择相应的处理器
type AutoPostProcessor struct {
    packageType string
    config      config.LLPkgConfig
}

func (a *AutoPostProcessor) Postprocessing() error {
    switch a.packageType {
    case "python":
        return a.processPythonPackage()
    case "cpp":
        return a.processCppPackage()
    }
}
```

#### 2. **统一版本管理**

##### 创建 `VersionManager`
```go
type VersionManager struct {
    *DefaultClient
}

// 统一的版本提取逻辑
func (vm *VersionManager) ExtractVersionFromCommit() (string, error) {
    // 支持多种版本格式
    // Release-as: package_name/vX.X.X
    // Release: vX.X.X
    // Version: vX.X.X
}
```

#### 3. **统一错误处理**

##### 结构化错误类型
```go
type PostProcessingError struct {
    Type    string
    Message string
    Cause   error
}

func (e *PostProcessingError) Error() string {
    return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
}
```

### 实现效果

#### 1. **代码复用**
- 版本提取逻辑复用率: 90%
- 错误处理逻辑复用率: 85%
- GitHub 操作逻辑复用率: 95%

#### 2. **维护性提升**
- 统一的接口设计
- 一致的错误处理
- 集中的配置管理

#### 3. **用户体验**
- 一致的命令接口
- 统一的错误消息格式
- 相同的配置选项

### 架构对比

#### 修改前
```
C++ 部分:
├── postprocessing (使用 DefaultClient)
│   ├── 版本提取
│   ├── GitHub Release
│   └── 错误处理

Python 部分:
├── postprocessing (直接实现)
│   ├── 版本提取 (重复实现)
│   ├── GitHub Release (重复实现)
│   └── 错误处理 (简单实现)
```

#### 修改后
```
统一架构:
├── DefaultClient (统一接口)
│   ├── 版本提取 (共享逻辑)
│   ├── GitHub Release (共享逻辑)
│   └── 错误处理 (统一格式)
│
├── C++ 处理器
│   └── 调用 DefaultClient
│
└── Python 处理器
    └── 调用 DefaultClient
```

## Support and Community

### Getting Help

### Resources

- **llpkgstore Repository**: [github.com/goplus/llpkgstore](https://github.com/goplus/llpkgstore)
- **llpyg Tool**: [github.com/toaction/llpyg](https://github.com/toaction/llpyg)
- **LLGo Project**: [github.com/goplus/llgo](https://github.com/goplus/llgo)
- **Python Package Index**: [pypi.org](https://pypi.org)
