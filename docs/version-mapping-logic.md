# llpkgstore 版本映射逻辑详解

## 概述

`llpkgstore` 采用统一的版本管理机制，为 Python 和 C++ 包提供一致的版本映射、记录和管理功能。本文档详细说明了版本映射的完整流程和实现逻辑。

## 核心设计原则

### 1. 统一版本管理
- **Python 和 C++ 包使用相同的版本管理逻辑**
- **共享 `DefaultClient` 接口和实现**
- **一致的版本提取、记录和验证机制**

### 2. 基于 Commit 消息的版本提取
- 从 Git commit 消息中提取版本信息
- 支持 `Release-as: {包名}/v{版本号}` 格式
- 自动解析包名和版本号

### 3. 集中式版本记录
- 使用 `llpkgstore.json` 文件记录版本映射
- 支持本地和集中式版本记录
- 兼容 C++ 和 Python 包的版本信息

## 版本映射流程

### 1. 版本提取阶段

#### 1.1 Commit 消息解析
```go
// 从 commit 消息中提取版本信息
// 示例: "Release-as: tabulate/v0.0.26"
func (d *DefaultClient) mappedVersion() (string, error) {
    sha, err := env.LatestCommitSHA()
    if err != nil {
        return "", err
    }
    
    // 检查是否与 PR 关联
    if !d.isAssociatedWithPullRequest(sha) {
        return "", fmt.Errorf("actions: not a merge request commit")
    }
    
    // 从 commit 消息中提取版本
    version, err := d.extractVersionFromCommit(sha)
    return version, err
}
```

#### 1.2 版本格式验证
支持的版本格式：
- `Release-as: {包名}/v{版本号}` (推荐)
- `Release-as: {包名}/{版本号}`
- 语义化版本号 (SemVer): `v1.2.3`, `1.2.3`

### 2. 配置文件处理阶段

#### 2.1 包名提取
```go
// 从版本字符串中提取包名和映射版本
func parseMappedVersion(version string) (clib, mappedVersion string, err error) {
    // 解析 "tabulate/v0.0.26" -> clib="tabulate", mappedVersion="v0.0.26"
    parts := strings.Split(version, "/")
    if len(parts) != 2 {
        return "", "", fmt.Errorf("invalid version format")
    }
    return parts[0], parts[1], nil
}
```

#### 2.2 配置文件读取
```go
// 读取包配置文件
cfg, err := config.ParseLLPkgConfig(filepath.Join(clib, "llpkg.cfg"))
if err != nil {
    return err
}
```

配置文件结构 (`llpkg.cfg`):
```json
{
  "type": "python",
  "upstream": {
    "installer": {
      "name": "pip"
    },
    "package": {
      "name": "tabulate",
      "version": "0.9.0"
    }
  },
  "llpyg": {
    "output_dir": "./test",
    "mod_name": "github.com/PengPengPeng717/llpkg/tabulate",
    "mod_depth": 1
  }
}
```

### 3. 版本记录更新阶段

#### 3.1 版本信息结构
```go
// LLPkgStoreJSON 表示版本记录文件结构
type LLPkgStoreJSON struct {
    Packages map[string]PackageInfo `json:"packages"`
}

// PackageInfo 表示包的版本信息
type PackageInfo struct {
    Versions []VersionInfo `json:"versions"`
}

// VersionInfo 表示单个版本信息
type VersionInfo struct {
    UpstreamVersion string `json:"upstream_version"`  // 上游版本 (如: "0.9.0")
    MappedVersion   string `json:"mapped_version"`    // 映射版本 (如: "v0.0.26")
    CreatedAt       string `json:"created_at"`        // 创建时间
}
```

#### 3.2 版本记录更新逻辑
```go
// 更新版本记录
func (v *Versions) Write(clib, upstreamVersion, mappedVersion string) {
    if v.Packages == nil {
        v.Packages = make(map[string]PackageInfo)
    }
    
    if v.Packages[clib].Versions == nil {
        v.Packages[clib] = PackageInfo{
            Versions: []VersionInfo{},
        }
    }
    
    // 添加新版本记录
    newVersion := VersionInfo{
        UpstreamVersion: upstreamVersion,
        MappedVersion:   mappedVersion,
        CreatedAt:       time.Now().Format(time.RFC3339),
    }
    
    v.Packages[clib].Versions = append(v.Packages[clib].Versions, newVersion)
}
```

### 4. Git 标签管理阶段

#### 4.1 标签创建
```go
// 创建 Git 标签
func (d *DefaultClient) createTag(version, sha string) error {
    // 检查标签是否已存在
    if hasTag(version) {
        return fmt.Errorf("actions: tag has already existed")
    }
    
    // 创建轻量级标签
    cmd := exec.Command("git", "tag", version, sha)
    return cmd.Run()
}
```

#### 4.2 标签推送
```go
// 推送标签到远程仓库
func pushTagToRemote(version string) error {
    cmd := exec.Command("git", "push", "origin", version)
    return cmd.Run()
}
```

### 5. GitHub Release 管理阶段

#### 5.1 Release 创建
```go
// 创建 GitHub Release
func (d *DefaultClient) createReleaseByTag(version string) (*github.RepositoryRelease, error) {
    release := &github.RepositoryRelease{
        TagName:    &version,
        Name:       &version,
        Body:       github.String("Release " + version),
        Draft:      github.Bool(false),
        Prerelease: github.Bool(false),
    }
    
    return d.client.Repositories.CreateRelease(context.Background(), d.owner, d.repo, release)
}
```

#### 5.2 Artifacts 上传
```go
// 上传构建产物到 Release
func (d *DefaultClient) uploadArtifactsToRelease(release *github.RepositoryRelease) ([]*os.File, error) {
    // 获取当前 workflow run 的 artifacts
    artifacts, _, err := d.client.Actions.ListWorkflowRunArtifacts(ctx, d.owner, d.repo, id, &github.ListOptions{})
    
    if artifacts.GetTotalCount() == 0 {
        // 没有 artifacts 时返回空结果，允许 postprocessing 继续
        return []*os.File{}, nil
    }
    
    // 上传每个 artifact 到 release
    for _, artifact := range artifacts.Artifacts {
        // 下载并上传 artifact
        err := d.uploadArtifactToRelease(artifact, release)
        if err != nil {
            return nil, err
        }
    }
    
    return files, nil
}
```

## 版本映射示例

### 示例 1: Python 包版本映射

**输入:**
- Commit 消息: `Release-as: tabulate/v0.0.26`
- 配置文件: `tabulate/llpkg.cfg`
- 上游版本: `0.9.0`

**处理流程:**
1. 提取包名: `tabulate`
2. 提取映射版本: `v0.0.26`
3. 读取配置文件获取上游版本: `0.9.0`
4. 更新版本记录:
   ```json
   {
     "packages": {
       "tabulate": {
         "versions": [
           {
             "upstream_version": "0.9.0",
             "mapped_version": "v0.0.26",
             "created_at": "2025-09-05T09:31:38Z"
           }
         ]
       }
     }
   }
   ```
5. 创建 Git 标签: `tabulate/v0.0.26`
6. 创建 GitHub Release: `tabulate/v0.0.26`

### 示例 2: C++ 包版本映射

**输入:**
- Commit 消息: `Release-as: clib/v1.0.0`
- 配置文件: `clib/llpkg.cfg`
- 上游版本: `1.0.0`

**处理流程:**
1. 提取包名: `clib`
2. 提取映射版本: `v1.0.0`
3. 读取配置文件获取上游版本: `1.0.0`
4. 更新版本记录
5. 创建 Git 标签: `clib/v1.0.0`
6. 创建 GitHub Release: `clib/v1.0.0`

## 版本记录文件管理

### 1. 本地版本记录
- 文件位置: `{包目录}/llpkgstore.json`
- 用途: 记录单个包的版本历史
- 格式: 与集中式版本记录相同

### 2. 集中式版本记录
- 文件位置: `llpkg/public/llpkgstore.json`
- 用途: 记录所有包的版本映射
- 格式: 包含所有包的版本信息

### 3. 版本记录更新策略
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

## 错误处理和验证

### 1. 版本格式验证
```go
// 验证版本格式
func isValidVersionFormat(version string) bool {
    // 支持 v1.2.3, 1.2.3 等格式
    versionRegex := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
    return versionRegex.MatchString(version)
}
```

### 2. 重复版本检查
```go
// 检查标签是否已存在
func hasTag(version string) bool {
    cmd := exec.Command("git", "tag", "-l", version)
    output, err := cmd.Output()
    if err != nil {
        return false
    }
    return strings.TrimSpace(string(output)) == version
}
```

### 3. 配置文件验证
```go
// 验证配置文件存在性和格式
func validateConfigFile(configPath string) error {
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return fmt.Errorf("config file not found: %s", configPath)
    }
    
    _, err := config.ParseLLPkgConfig(configPath)
    return err
}
```

## 环境变量和配置

### 必需的 GitHub Actions 环境变量
- `GITHUB_TOKEN`: GitHub API 访问令牌
- `GITHUB_REPOSITORY`: 仓库名称 (格式: `owner/repo`)
- `GITHUB_SHA`: 当前 commit SHA
- `GITHUB_REF`: 当前分支或标签引用

### 可选环境变量
- `GITHUB_ACTOR`: 触发 workflow 的用户
- `GITHUB_EVENT_NAME`: 触发事件名称
- `GITHUB_WORKFLOW`: workflow 名称

## 最佳实践

### 1. Commit 消息规范
- 使用 `Release-as: {包名}/v{版本号}` 格式
- 版本号遵循语义化版本规范
- 确保包名与目录名一致

### 2. 配置文件管理
- 确保 `llpkg.cfg` 文件格式正确
- 定期验证配置文件的有效性
- 保持上游版本信息的准确性

### 3. 版本记录维护
- 定期备份版本记录文件
- 监控版本映射的一致性
- 及时处理版本冲突

### 4. 错误处理
- 实现适当的重试机制
- 记录详细的错误日志
- 提供清晰的错误信息

## 总结

`llpkgstore` 的版本映射逻辑采用统一的设计，为 Python 和 C++ 包提供一致的版本管理体验。通过基于 commit 消息的版本提取、集中式版本记录和自动化 Git 标签管理，确保了版本映射的准确性和可追溯性。

关键特性：
- ✅ **统一版本管理**: Python 和 C++ 包使用相同的逻辑
- ✅ **自动化版本提取**: 从 commit 消息中自动提取版本信息
- ✅ **集中式版本记录**: 统一的版本映射记录和管理
- ✅ **Git 标签管理**: 自动创建和推送版本标签
- ✅ **GitHub Release 集成**: 自动创建 Release 并上传构建产物
- ✅ **错误处理和验证**: 完善的错误处理和版本验证机制
