# llpyg 配置选项

`llpkgstore` 现在支持在 `llpkg.cfg` 文件中配置 `llpyg` 的命令行参数。

## 配置结构

在 `llpkg.cfg` 文件中添加 `llpyg` 配置段：

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
    "output_dir": "./generated",
    "mod_name": "github.com/PengPengPeng717/llpkg/numpy",
    "mod_depth": 2
  }
}
```

## 配置选项说明

### `output_dir` (可选)
- **类型**: `string`
- **默认值**: `"./test"`
- **说明**: 指定 `llpyg` 的输出目录，对应 `llpyg` 的 `-o` 参数

### `mod_name` (可选)
- **类型**: `string`
- **默认值**: 包名
- **说明**: 指定生成的 Go 模块名，对应 `llpyg` 的 `-mod` 参数
- **示例**: `"github.com/PengPengPeng717/llpkg/numpy"`

### `mod_depth` (可选)
- **类型**: `int`
- **默认值**: `1`
- **范围**: `0-10`
- **说明**: 指定 Python 模块的最大提取深度，对应 `llpyg` 的 `-d` 参数

## 使用示例

### 基本配置
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
使用默认参数：输出目录 `./test`，模块名 `numpy`，深度 `1`

### 完整配置
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
    "mod_name": "github.com/myorg/mypackage/numpy",
    "mod_depth": 3
  }
}
```

## 执行命令

配置完成后，执行：
```bash
llpkgstore generate
```

这相当于执行：
```bash
llpyg -o ./bindings -mod github.com/myorg/mypackage/numpy -d 3 numpy
```

## 验证规则

- `mod_depth` 必须为非负数且不超过 10
- 如果配置无效，`llpkgstore generate` 会报错并退出 