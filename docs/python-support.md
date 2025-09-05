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

### Unified Version Management (v2.0+)

Python packages now use the same version management mechanism as C/C++ packages, providing a consistent experience across all package types.

### Commit Message Format

Use this format for version releases:

```bash
git commit -m "Release-as: {package_name}/v{version}"
```

Examples:
```bash
git commit -m "Release-as: numpy/v1.26.4"
git commit -m "Release-as: requests/v2.31.0"
git commit -m "Release-as: tabulate/v10.0.0"
```

### Smart Version Extraction

The system automatically extracts version information from:
1. **Latest commit messages** (priority for CI environments)
2. **Git tags** (fallback when commit messages don't contain version info)

Supported formats:
- `Release-as: package_name/vX.X.X`
- `Release: vX.X.X`
- `Version: vX.X.X`

### Dual Version Recording

The system maintains version information in two locations:

#### Local Version Record
File: `{package_dir}/llpkgstore.json`
```json
{
  "packages": {
    "tabulate": {
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

#### Centralized Version Record
File: `llpkg/public/llpkgstore.json`
```json
{
  "packages": {
    "tabulate": {
      "versions": [
        {
          "python": "0.9.0",
          "go": ["v0.0.2", "v10.0.0"]
        }
      ]
    }
  }
}
```

### Automatic Git Tagging

The system automatically:
- **Creates Git tags** based on extracted version information
- **Pushes tags** to remote repository
- **Checks for duplicates** to avoid tag conflicts
- **Provides detailed logging** for tag operations

### Postprocessing Command

Use the postprocessing command to handle version management:

```bash
llpkgstore postprocessing
```

This command will:
1. Extract version from commit messages or Git tags
2. Update local `llpkgstore.json` with version mapping
3. Update centralized `llpkg/public/llpkgstore.json`
4. Create and push Git tags
5. Create GitHub releases (in CI environment)

### Version Compatibility

- Python versions are mapped to Go versions using unified format
- Follows semantic versioning standards
- Maintains backward compatibility
- Supports multiple version branches
- Provides conflict detection and resolution

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

## Support and Community

### Getting Help

- **GitHub Issues**: Report bugs and request features
- **Documentation**: Comprehensive guides and examples
- **Community**: Active developer community
- **Examples**: Rich collection of working examples

### Contributing

- **Bug Reports**: Detailed issue descriptions
- **Feature Requests**: Clear use case descriptions
- **Code Contributions**: Follow contribution guidelines
- **Documentation**: Help improve documentation

### Resources

- **llpkgstore Repository**: [github.com/goplus/llpkgstore](https://github.com/goplus/llpkgstore)
- **llpyg Tool**: [github.com/toaction/llpyg](https://github.com/toaction/llpyg)
- **LLGo Project**: [github.com/goplus/llgo](https://github.com/goplus/llgo)
- **Python Package Index**: [pypi.org](https://pypi.org)

---

For more information, see the main [llpkgstore documentation](./llpkgstore.md) or visit our [GitHub repository](https://github.com/goplus/llpkgstore).
