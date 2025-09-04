# llpkgstore

[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/llpkgstore)](https://goreportcard.com/report/github.com/goplus/llpkgstore)
[![Go Version](https://img.shields.io/github/go-mod/go-version/goplus/llpkgstore)](https://github.com/goplus/llpkgstore/blob/main/go.mod)
[![License](https://img.shields.io/github/license/goplus/llpkgstore)](https://github.com/goplus/llpkgstore/blob/main/LICENSE)

**llpkgstore** is a comprehensive package distribution service for [**LLGo**](https://github.com/goplus/llgo), designed to provide trustworthy and convenient access to language bindings for multiple ecosystems.

## 🚀 Features

### ✨ Multi-Language Support
- **C/C++ Packages**: Full support via `llcppg` tool
- **Python Packages**: Complete support via `llpyg` tool (v2.0+)

### 🔧 Core Capabilities
- **Automated Generation**: CI/CD pipeline for package creation
- **Version Management**: Intelligent version mapping and tagging
- **Type Safety**: Maintained type information across language boundaries
- **Testing Framework**: Automated demo code generation and testing
- **GitHub Integration**: Seamless GitHub Actions and Releases support

### 🎯 Use Cases
- **Library Integration**: Use C/C++ libraries in Go
- **Python Ecosystem**: Access Python packages from Go
- **Cross-Platform**: Consistent experience across different platforms
- **Enterprise Ready**: Production-grade package management

## 📦 Supported Package Types

### C/C++ Packages
- **Tool**: `llcppg`
- **Installer**: Conan package manager
- **Status**: ✅ Fully supported
- **Features**: Binary distribution, header files, `.pc` files

### Python Packages
- **Tool**: `llpyg`
- **Installer**: pip package manager
- **Status**: ✅ Fully supported (v2.0+)
- **Features**: Python module bindings, Go interfaces, type safety

## 🚀 Quick Start

### 1. Install llpkgstore

```bash
go install github.com/goplus/llpkgstore/cmd/llpkgstore@latest
```

### 2. Create a Python Package

```bash
mkdir my-python-package
cd my-python-package
```

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

### 4. Test the Package

```bash
cd _demo/basic_test
llgo run main.go
```

### 5. Use in Your Go Code

```go
package main

import (
    "github.com/goplus/lib/py"
    "requests"
)

func main() {
    // Use the generated Python package
    resp := requests.Get("https://api.github.com")
    fmt.Println(resp.StatusCode)
}
```

## 📚 Documentation

- **[Main Documentation](./docs/llpkgstore.md)**: Comprehensive design and architecture guide
- **[Python Support](./docs/python-support.md)**: Detailed Python package support guide
- **[llpyg Configuration](./docs/llpyg-config.md)**: Configuration options for Python packages

## 🏗️ Architecture

llpkgstore is built with a modular architecture that supports multiple package types:

```
llpkgstore/
├── cmd/llpkgstore/           # CLI application
│   ├── internal_python/      # Python package support
│   └── internal_cpp/         # C/C++ package support
├── internal/                  # Core functionality
│   ├── actions/              # GitHub Actions integration
│   ├── config/               # Configuration management
│   └── versions/             # Version management
├── upstream/                  # Package installer support
│   ├── pip/                  # Python package installer
│   └── conan/                # C/C++ package installer
└── docs/                     # Documentation
```

## 🔄 Workflow

### Package Generation Flow

1. **Configuration**: Define package in `llpkg.cfg`
2. **Generation**: Run `llpkgstore generate`
3. **Testing**: Execute demo code for verification
4. **Commit**: Use proper version format in commit messages
5. **CI/CD**: Automated testing and release via GitHub Actions
6. **Distribution**: Package available via `llgo get`

### Version Management

- **Commit Format**: `Release-as: {package_name}/v{version}`
- **Automatic Tagging**: GitHub releases and tags
- **Version Mapping**: Intelligent version conversion
- **Branch Management**: Automated branch lifecycle

## 🛠️ Development

### Prerequisites

- Go 1.24+
- Python 3.8+ (for Python packages)
- LLGo toolchain
- llpyg tool (for Python packages)

### Building from Source

```bash
git clone https://github.com/goplus/llpkgstore.git
cd llpkgstore
go build -o llpkgstore ./cmd/llpkgstore
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📊 Status

### Current Status
- **C/C++ Support**: ✅ Production Ready
- **Python Support**: ✅ Production Ready (v2.0+)
- **CI/CD Pipeline**: ✅ Fully Automated
- **Documentation**: ✅ Comprehensive

### Roadmap
- **Rust Support**: 🚧 In Development
- **Node.js Support**: 📋 Planned
- **Java Support**: 📋 Planned
- **Performance Optimization**: 🔄 Ongoing

## 🤝 Community

### Getting Help
- **GitHub Issues**: [Report bugs and request features](https://github.com/goplus/llpkgstore/issues)
- **Discussions**: [Join community discussions](https://github.com/goplus/llpkgstore/discussions)
- **Documentation**: [Comprehensive guides](./docs/)

### Contributing
- **Code**: Follow Go coding standards
- **Documentation**: Help improve guides and examples
- **Testing**: Add tests for new features
- **Examples**: Share use cases and examples

### Resources
- **LLGo Project**: [github.com/goplus/llgo](https://github.com/goplus/llgo)
- **llpyg Tool**: [github.com/toaction/llpyg](https://github.com/toaction/llpyg)
- **llcppg Tool**: [github.com/goplus/llcppg](https://github.com/goplus/llcppg)

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **LLGo Team**: For the amazing LLGo ecosystem
- **llpyg Contributors**: For Python binding support
- **llcppg Contributors**: For C/C++ binding support
- **Community**: For feedback and contributions

---

**Made with ❤️ by the LLGo community**

For more information, visit [llpkg.goplus.org](https://llpkg.goplus.org) or check out our [GitHub repository](https://github.com/goplus/llpkgstore).
