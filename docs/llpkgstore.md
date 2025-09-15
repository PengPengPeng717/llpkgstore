# llpkgstore 技术设计文档

本文档详细描述了 **llpkgstore** 的技术设计、架构实现和核心机制。

## 概述

llpkgstore 是一个专为 [**LLGo**](https://github.com/goplus/llgo) 设计的统一包分发服务，为多语言生态系统提供可信赖且便捷的语言绑定访问。

### 设计目标

llpkgstore 旨在解决以下问题：

1. **跨语言集成复杂性**: 简化不同编程语言之间的互操作
2. **包管理分散**: 统一管理多种语言的包和绑定
3. **版本管理混乱**: 提供一致的版本映射和管理机制
4. **安全性问题**: 通过自动化流程确保包的可信度

### 核心概念

**llpkg** 是一个 Go 模块，通过 [**LLGo**](https://github.com/goplus/llgo) 的生态系统集成能力调用其他语言的库。目前，llpkg 生成由以下工具处理：

- [**`llcppg`**](https://github.com/goplus/llcppg) - C/C++ 库绑定生成器
- [**`llpyg`**](https://github.com/toaction/llpyg) - Python 包绑定生成器

### 系统组件

llpkgstore 由以下核心组件构成：

1. **CLI 工具**: 统一的命令行接口，提供包管理功能
2. **配置系统**: 基于 JSON 的配置文件格式，支持多种包类型
3. **版本管理**: 智能版本提取、映射和记录机制
4. **生成器集成**: 与 llcppg 和 llpyg 的无缝集成
5. **CI/CD 流水线**: 自动化包生成、测试和发布流程

## 支持的包类型

### C/C++ 包

**工具**: `llcppg`
**包管理器**: Conan
**状态**: ✅ 完全支持
**特性**:
- 二进制分发支持
- 头文件自动处理
- `.pc` 文件生成
- 跨平台兼容性

**配置示例**:
```json
{
  "type": "cpp",
  "upstream": {
    "installer": {
      "name": "conan"
    },
    "package": {
      "name": "opencv",
      "version": "4.8.0"
    }
  }
}
```

### Python 包

**工具**: `llpyg`
**包管理器**: pip
**状态**: ✅ 完全支持 (v2.0+)
**特性**:
- Python 模块绑定
- Go 接口生成
- 类型安全保证
- 智能包检测

**配置示例**:
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
  },
  "llpyg": {
    "output_dir": "./bindings",
    "mod_name": "github.com/your-org/numpy",
    "mod_depth": 1
  }
}
```

### 未来支持

- **Rust**: 计划在未来版本中支持
- **Node.js**: 正在考虑中
- **Java**: 正在考虑中

## 目录结构

### C/C++ 包结构

```
{CLibraryName}/
├── llpkg.cfg              # 包配置文件（必需）
├── llcppg.cfg             # llcppg 配置文件
├── llcppg.symb.json       # 符号定义文件
├── llcppg.pub             # 公共接口文件
├── go.mod                 # Go 模块文件
├── go.sum                 # Go 依赖校验文件
├── {GeneratedGoFiles}     # 生成的 Go 绑定文件
└── _demo/                 # 演示和测试代码
    ├── {DemoName1}/
    │   └── main.go
    └── {DemoName2}/
        └── main.go
```

### Python 包结构

```
{PythonPackageName}/
├── llpkg.cfg              # 包配置文件（必需）
├── llpyg.cfg              # llpyg 配置文件
├── go.mod                 # Go 模块文件
├── go.sum                 # Go 依赖校验文件
├── {GeneratedGoFiles}     # 生成的 Go 绑定文件
│   ├── {PackageName}.go
│   └── {PackageName}_autogen_link.go
└── _demo/                 # 演示和测试代码
    ├── {DemoName1}/
    │   └── main.go
    └── {DemoName2}/
        └── main.go
```

### 文件说明

| 文件 | 描述 | 必需性 |
|------|------|--------|
| `llpkg.cfg` | 包配置文件，定义包类型和基本信息 | ✅ 必需 |
| `llcppg.cfg` | C/C++ 包生成器配置 | C/C++ 包必需 |
| `llcppg.symb.json` | C/C++ 符号定义文件 | C/C++ 包必需 |
| `llcppg.pub` | C/C++ 公共接口文件 | C/C++ 包必需 |
| `llpyg.cfg` | Python 包生成器配置 | Python 包必需 |
| `go.mod` | Go 模块定义文件 | 自动生成 |
| `go.sum` | Go 依赖校验文件 | 自动生成 |
| `_demo/` | 演示和测试代码目录 | 自动生成 |

### 配置文件要求

为了确保 `llgo` 能够正确识别 llpkg，每个包必须包含至少一个 `llpkg.cfg` 文件。该文件定义了包的基本信息，包括：

- 包类型（python 或 cpp）
- 上游包信息
- 生成器配置选项

## llpkg.cfg Structure

### C/C++ Package Configuration
```json
{
  "type": "cpp",
  "upstream": {
    "installer": {
      "name": "conan",
      "config": {
        "options": ""
      }
    },
    "package": {
      "name": "cjson",
      "version": "1.7.18"
    }
  }
}
```

### Python Package Configuration
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
      "name": "numpy",
      "version": "1.26.4"
    }
  },
  "llpyg": {
    "output_dir": "./test",
    "mod_name": "github.com/PengPengPeng717/llpkg/numpy",
    "mod_depth": 1
  }
}
```

### Field description

**Common Fields**

| key | type | defaultValue | optional | description |
|------|------|--------|------|------|
| type | `string` | "cpp" | ✅ | Package type: "cpp" or "python" |
| upstream.installer.name | `string` | "conan" | ✅ | Upstream binary provider |
| upstream.installer.config | `map[string]string` | {} | ✅ | Config of installer |
| upstream.package.name | `string` | - | ❌ | Package name in platform |
| upstream.package.version | `string` | - | ❌ | Original package version |

**Python-specific Fields (llpyg section)**

| key | type | defaultValue | optional | description |
|------|------|--------|------|------|
| llpyg.output_dir | `string` | "./test" | ✅ | Output directory for generated files |
| llpyg.mod_name | `string` | package name | ✅ | Go module name |
| llpyg.mod_depth | `int` | 1 | ✅ | Maximum module extraction depth (0-10) |

## Getting an llpkg

### C/C++ Packages
Use `llgo get` to get a C/C++ llpkg:

```bash
llgo get clib@cversion
```

*e.g.* `llgo get cjson@1.7.18`

- `clib`: the original library name in C
- `cversion`: the original version in C

### Python Packages
Use `llgo get` to get a Python llpkg:

```bash
llgo get github.com/goplus/llpkg/numpy@v1.26.4
```

Or use the simplified syntax (if supported):

```bash
llgo get numpy@1.26.4
```

### Universal Syntax
Both package types support the universal syntax:

```bash
llgo get module_path@module_version
```

*e.g.* `llgo get github.com/goplus/llpkg/cjson@v1.0.0`

```bash
llgo get clib[@latest]
llgo get module_path[@latest]
```

The optional `latest` identifier is supported as a valid `cversion` or `module_version`. When `llgo get clib@latest`, `llgo get` will firstly convert `clib` to `module_path`, and then process it as `module_path@latest`. `llgo get` will find the latest llpkg and pull it.

Wrong usage:

```bash
llgo get clib@module_version
llgo get module_path@cversion
```

It's the format of the part before `@` that determines the how `llgo get` will handle the version; that is, `llgo get` will firstly check if it's a `clib`. If it is, the whole argument will be processed as `clib@cversion`; otherwise, it will be processed as `module_path@module_version`.

> **Details of `llgo get`**
>
>  1. `llgo` automatically resolves `clib@cversion` syntax into canonical `module_path@module_version` format.
>  2. Pull the go module by `go get`.
>  3. Check `llpkg.cfg` to determine if it's an llpkg. If it is:
>    - For C/C++ packages: `llgo get` will run `upstream.installer` to install binaries. `.pc` files for building will be stored in `{LLGOPCCACHE}`.
>    - For Python packages: `llgo get` will use the generated Go bindings directly.
>    - A comment in `go.mod` will be added to indicate the original version. Comments of indirect dependencies will be automatically processed by `go mod tidy`.
>
>       ```
>       // go.mod for C/C++ package
>       require (
>             github.com/goplus/llpkg/cjson v1.1.0  // conan:cjson/1.7.18
>       )
>
>       // go.mod for Python package
>       require (
>             github.com/goplus/llpkg/numpy v1.26.4  // pip:numpy/1.26.4
>       )
>       ```

## Listing clib version mapping

```
llgo list -m [-versions] [-json] [modules/clibs]
```

- `llgo list -m` is compatible with `go list -m`.
- `modules`: a set of space-separated module_path[@module_version].
- `clibs`: a set of space-separated clib[@cversion]

Each argument is processed separately.

### `module`

`llgo list` will check if the `module` is an llpkg or a normal go module by seeking if `llpkg.cfg` exists.

#### llpkg

If the `module` is an llpkg:

1. `llgo list -m`

`llgo list` will print the module path and the upstream of the local llpkg according to `go.mod` and `llpkg.cfg`.

*e.g.* `llgo list -m cjson` (C/C++ package):

```
github.com/goplus/llpkg/cjson v0.1.0[conan:cjson/1.7.18]
```

*e.g.* `llgo list -m numpy` (Python package):

```
github.com/goplus/llpkg/numpy v1.26.4[pip:numpy/1.26.4]
```

2. `llgo list -m -versions`

Add `-versions` to check all version mappings of an llpkg.

*e.g.* `llgo list -m -versions cjson` or `llgo list -m -versions github.com/goplus/llpkg/cjson`:

```
github.com/goplus/llpkg/cjson v0.1.0[conan:cjson/1.7.18] v0.1.1[conan:cjson/1.7.18] v0.2.0[conan:cjson/1.7.19]
```

*e.g.* `llgo list -m -versions numpy`:

```
github.com/goplus/llpkg/numpy v1.26.4[pip:numpy/1.26.4] v1.27.0[pip:numpy/1.27.0]
```

3. JSON output

We define a Go Struct for the output of `llgo list -m -versions -json`:

```go
type Module struct {
	Path     string   `json:"Path"`
	Version  string   `json:"Version"`
	Versions []string `json:"Versions,omitempty"`
	Upstream string   `json:"Upstream,omitempty"`
}
```

*e.g.* `llgo list -m -versions -json cjson`:

```json
{
  "Path": "github.com/goplus/llpkg/cjson",
  "Version": "v0.2.0",
  "Versions": ["v0.1.0", "v0.1.1", "v0.2.0"],
  "Upstream": "conan:cjson/1.7.19"
}
```

*e.g.* `llgo list -m -versions -json numpy`:

```json
{
  "Path": "github.com/goplus/llpkg/numpy",
  "Version": "v1.27.0",
  "Versions": ["v1.26.4", "v1.27.0"],
  "Upstream": "pip:numpy/1.27.0"
}
```

### `clib`

`llgo list` will convert `clib` to `module_path` and then process it as `module`.

*e.g.* `llgo list -m cjson`:

```
github.com/goplus/llpkg/cjson v0.2.0[conan:cjson/1.7.19]
```

*e.g.* `llgo list -m numpy`:

```
github.com/goplus/llpkg/numpy v1.27.0[pip:numpy/1.27.0]
```

## Package Generation Workflow

### C/C++ Package Generation

A standard method for generating valid C/C++ llpkgs:
1. Receive binaries/headers from [installer](#llpkgcfg-structure), and index them into `.pc` files
2. Detect the generator from configuration files. For example, if an `llcppg.cfg` file is present in the current directory, we can directly use `llcppg`
3. Automatically generate llpkg using a generator for different platforms
4. Combine generated results into one Go module
5. Debug and re-generate llpkg by modifying the configuration file

### Python Package Generation

A standard method for generating valid Python llpkgs:
1. Install Python package using pip installer
2. Use `llpyg` tool to generate Go bindings for Python modules
3. Configure output directory, module name, and extraction depth
4. Generate Go interfaces and type-safe bindings
5. Create Go module with proper dependencies
6. Test generated bindings with demo code

### Generation Commands

#### C/C++ Packages
```bash
llpkgstore generate
```

#### Python Packages
```bash
llpkgstore generate
```

This automatically detects the package type from `llpkg.cfg` and uses the appropriate generator.

## PR Workflow

### Standard PR workflow
1. Create PR to trigger GitHub Action
2. PR verification
3. llpkg generation
4. Run test
5. Review generated llpkg
6. Merge PR
7. Run post-processing Github Action on main branch

### PR verification workflow
1. Ensure that there is only one `llpkg.cfg` file across all directories. If multiple instances of `llpkg.cfg` are detected, the PR will be aborted.
2. Check if the directory name is valid, the directory name in PR **SHOULD** equal to `Package.Name` field in the `llpkg.cfg` file.
3. Check the PR commit footer contains a [`{MappedVersion}`](#mappedversion-in-pr-commit).

### llpkg generation

A standard method for generating valid llpkgs:
1. **C/C++ packages**: Receive binaries/headers from [installer](#llpkgcfg-structure), and index them into `.pc` files
2. **Python packages**: Install Python package using pip and generate Go bindings with `llpyg`
3. Detect the generator from configuration files. For example, if an `llcppg.cfg` file is present in the current directory, we can directly use `llcppg`, or if `llpyg.cfg` is present, we can use `llpyg`
4. Automatically generate llpkg using a generator for different platforms
5. Combine generated results into one Go module
6. Debug and re-generate llpkg by modifying the configuration file

### Merge PR
The maintainer **SHOULD** squash commits before merging a PR. The squash commit message **MUST** include [`{MappedVersion}`](#mappedversion-in-pr-commit) to enable the Post-processing GitHub Action to parse it correctly.

#### `{MappedVersion}` in PR Commit
The `{MappedVersion}` **MUST** be included in at least one of the commits in the PR and **MUST** follow this format:

```
Release-as: {PackageName}/{MappedVersion}
```

The PR verification process will validate this format and abort the PR if it is invalid.

**Example for C/C++ package:**
```bash
git merge
# Modify the merge commit message
git commit --amend -m "feat: add cjson" -m "Release-as: cjson/v1.0.0"
```

**Example for Python package:**
```bash
git merge
# Modify the merge commit message
git commit --amend -m "feat: add numpy" -m "Release-as: numpy/v1.26.4"
```

### Post-processing GitHub Action
The Post-processing GitHub Action will tag the commit according to the [Version Tag Rule](#version-tag-rule).

#### Version Tag Rule
1. Extract the `{MappedVersion}` of the current package from the footer of the squashed commit.
2. Follow Go's version management for nested modules and tag `{PackageName}/{MappedVersion}` for each version.
3. This design is fully compatible with native Go modules:
    ```
    github.com/goplus/llpkg/cjson@v1.7.18
    github.com/goplus/llpkg/numpy@v1.26.4
    ```

### Legacy version maintenance workflow

1. Create an issue to discuss the package that requires maintenance.
2. The maintainer creates a label in the format `branch:release-branch.{PackageName}/{MappedVersion}` and adds it to the issue if the package needs maintenance.
3. A GitHub Action is triggered when the label is created. It determines whether a branch should be created based on the [Branch Maintenance Strategy](#branch-maintenance-strategy).
4. Open a pull request (PR) for maintenance. The maintainer **SHOULD** merge the PR with the commit message `fixed {IssueID}` to close the related issue.
5. When issues labeled with `branch:release-branch.` are closed, we need to determine whether to remove the branch. In the following case, the branch and label can be safely removed:
   - No associated PR with commit containing `fix* {ThisIssueID}`.(* means the commit starting with `fix` prefix)

## llpkg.goplus.org

This service is hosted by GitHub Pages, and the `llpkgstore.json` file is located in the same branch as GitHub Pages. When running `llgo get`, it will download the file to `LLGOPCCACHE`.

### Function

1. Provide a download of the mapping table.
2. Provide version queries for Go Modules corresponding to C libraries and Python packages.
3. Provide links to specific C libraries on Conan.io and Python packages on PyPI.

### Router

1. `/`: Home page with a search bar at the top and multiple llpkgs. Users can search for llpkgs by name and view the latest two versions. Clicking an llpkg opens a modal displaying:
   - Information about the original library (C library on Conan or Python package on PyPI)
   - All available versions of the llpkg

  ![Index](./llpkg_index.svg)

  ![Pkg detail](./llpkg_pkg.svg)

2. `/llpkgstore.json`: Provides the mapping table download.

**Note**: llpkg details are displayed in modals instead of new pages, as `llpkgstore.json` is loaded during the initial homepage access and does not require additional requests.

### Interaction with web service

When executing `llgo get clib@cversion` or `llgo get python_package@version`, a series of actions will be performed to map the version to `module_version`:
1. Fetch the latest `llpkgstore.json`
2. Parse the JSON file to find the corresponding `module_version` array
3. Select the latest patched version from the array
4. Retrieve llpkg

## Environment variable design

One usage is to store `.pc` files of the C library and allow `llgo build` to find them.

1. `LLGOCACHE` defaults to `{UserCacheDir}/llgo/`
2. `.pc` files of C libs needed by llpkg will be stored in `{LLGOCACHE}/pkg-config/{module_path}@{module_version}/`
3. If `UserCacheDir` isn't avaliable, `llgo` will exit with an error

## Python Package Support Details

### Python Package Generation

Python packages are generated using the `llpyg` tool, which creates Go bindings for Python modules. The process includes:

1. **Package Installation**: Uses pip to install the specified Python package
2. **Binding Generation**: Generates Go interfaces and type-safe bindings
3. **Module Configuration**: Configures output directory, module name, and extraction depth
4. **Go Module Creation**: Creates a proper Go module with dependencies

### Python Package Features

- **Type Safety**: Generated Go code includes proper type information
- **Module Depth Control**: Configurable extraction depth for nested modules
- **Custom Module Names**: Support for custom Go module names
- **Demo Code**: Automatic generation of test and demo code

### Python Package Testing

Python packages include demo code in the `_demo` directory to verify:
- Package import and compilation
- Basic functionality testing
- Type safety verification
- Integration with Go ecosystem

### Python Package Version Management

Python packages follow the same version management strategy as C/C++ packages:
- Version extraction from commit messages
- GitHub Release creation
- Tag management
- Version mapping in `llpkgstore.json`

### Smart Package Detection

llpkgstore v2.0+ includes intelligent package detection that prioritizes system-installed packages:

```go
// isPackageInstalledInSystem checks if the specified package is already installed in the system environment
func isPackageInstalledInSystem(packageName string) bool {
    // Method 1: Try to import the package directly
    if canImportPackage(packageName) {
        return true
    }
    
    // Method 2: Check pip list output
    if isPackageInPipList(packageName) {
        return true
    }
    
    return false
}
```

This optimization:
- **Reduces installation time**: From 6 minutes to almost 0 seconds
- **Saves network usage**: Avoids redundant downloads
- **Improves user experience**: Faster, smarter workflow
- **Maintains compatibility**: Falls back to temporary installation when needed

### llpyg Configuration

The `llpyg` section in `llpkg.cfg` provides fine-grained control over Python package generation:

```json
{
  "llpyg": {
    "output_dir": "./test",
    "mod_name": "github.com/PengPengPeng717/llpkg/numpy",
    "mod_depth": 1
  }
}
```

#### Configuration Options

| Field | Type | Default | Optional | Description |
|-------|------|---------|----------|-------------|
| `output_dir` | `string` | `"./test"` | ✅ | Output directory for generated files |
| `mod_name` | `string` | package name | ✅ | Go module name |
| `mod_depth` | `int` | `1` | ✅ | Maximum module extraction depth (0-10) |

#### Validation Rules

- `mod_depth` must be non-negative and not exceed 10
- `output_dir` cannot contain illegal characters (`..`, `~`, `\`, `:`, `*`, `?`, `"`, `<`, `>`, `|`)
- `mod_name` must contain at least one slash and start with a valid domain

## Architecture Improvements

### Unified Version Management

The latest version of llpkgstore includes unified version management for all package types:

1. **Common Version Extraction**: Unified logic for extracting versions from commit messages
2. **Standardized Version Validation**: Consistent semver validation across package types
3. **Unified Version Mapping**: Common version mapping strategies
4. **Consistent Branch Management**: Standardized branch naming and lifecycle

### Enhanced Error Handling

- **Structured Error Types**: Consistent error handling across all operations
- **Detailed Logging**: Comprehensive logging for debugging and monitoring
- **Error Recovery**: Mechanisms for handling and recovering from errors

### Improved CI/CD Pipeline

- **Package Type Detection**: Automatic detection of package types
- **Conditional Processing**: Different processing logic for different package types
- **Unified Workflow**: Consistent workflow across all package types
- **Enhanced Testing**: Comprehensive testing for all package types

## Future Roadmap

### Planned Features

1. **Rust Package Support**: Integration with Rust ecosystem
2. **Node.js Package Support**: JavaScript/TypeScript package bindings
3. **Java Package Support**: JVM ecosystem integration
4. **Enhanced Version Management**: Advanced versioning strategies
5. **Package Dependencies**: Cross-package dependency management

### Performance Improvements

1. **Parallel Processing**: Concurrent package generation
2. **Caching Mechanisms**: Intelligent caching for faster builds
3. **Incremental Generation**: Delta updates for existing packages
4. **Resource Optimization**: Better memory and CPU utilization

### Developer Experience

1. **IDE Integration**: Better editor support
2. **Debugging Tools**: Enhanced debugging capabilities
3. **Documentation**: Comprehensive API documentation
4. **Examples**: Rich examples and tutorials
