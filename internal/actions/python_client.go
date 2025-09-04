package actions

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goplus/llpkgstore/config"
	"github.com/goplus/llpkgstore/internal/actions/env"
	"github.com/goplus/llpkgstore/internal/actions/versions"
)

// PythonPostProcessor extends DefaultClient for Python-specific operations
type PythonPostProcessor struct {
	*DefaultClient
}

// NewPythonPostProcessor creates a new Python post-processor
func NewPythonPostProcessor() (*PythonPostProcessor, error) {
	client, err := NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create default client: %v", err)
	}

	return &PythonPostProcessor{
		DefaultClient: client,
	}, nil
}

// Postprocessing handles Python package post-processing using unified interface
// 完全遵循 C++ 的架构和流程
func (p *PythonPostProcessor) Postprocessing() error {
	fmt.Println("Starting unified Python package post-processing...")

	// 阶段 1: 版本检测与验证 (与 C++ 一致)
	sha, err := env.LatestCommitSHA()
	if err != nil {
		return err
	}

	// 检查是否为合并提交
	if !p.isAssociatedWithPullRequest(sha) {
		return fmt.Errorf("actions: not a merge request commit")
	}

	// 从提交消息中提取版本信息
	version, err := p.mappedVersion()
	if err != nil {
		return err
	}

	// 解析版本格式: "clib/semver"
	clib, mappedVersion, err := parseMappedVersion(version)
	if err != nil {
		return err
	}

	// 阶段 2: 配置解析 (与 C++ 一致)
	cfg, err := config.ParseLLPkgConfig(filepath.Join(clib, "llpkg.cfg"))
	if err != nil {
		return err
	}

	// 阶段 3: 版本映射存储 (与 C++ 一致)
	ver := versions.Read("llpkgstore.json")
	ver.Write(clib, cfg.Upstream.Package.Version, mappedVersion)

	// 阶段 4: Git 标签创建 (与 C++ 一致)
	if hasTag(version) {
		return fmt.Errorf("actions: tag has already existed")
	}

	if err := p.createTag(version, sha); err != nil {
		return err
	}

	// 阶段 5: GitHub Release 创建 (与 C++ 一致)
	release, err := p.createReleaseByTag(version)
	if err != nil {
		return err
	}

	_, err = p.uploadArtifactsToRelease(release)
	if err != nil {
		return err
	}

	// 阶段 6: 遗留分支清理 (与 C++ 一致)
	branchName, isLegacy, err := p.isLegacyVersion()
	if err != nil {
		return err
	}
	if isLegacy {
		err = p.removeBranch(branchName)
	}

	fmt.Println("Python package post-processing completed successfully")
	return err
}

// Python 特定的辅助方法可以在这里添加
// 目前使用与 C++ 完全一致的流程

// Python-specific helper methods

// ValidatePythonPackage validates Python package configuration
func (p *PythonPostProcessor) ValidatePythonPackage(cfg config.LLPkgConfig) error {
	if cfg.Type != "python" {
		return fmt.Errorf("package type must be 'python', got: %s", cfg.Type)
	}

	if cfg.Upstream.Installer.Name != "pip" {
		return fmt.Errorf("Python packages must use 'pip' installer, got: %s", cfg.Upstream.Installer.Name)
	}

	if cfg.Upstream.Package.Name == "" {
		return fmt.Errorf("package name is required")
	}

	if cfg.Upstream.Package.Version == "" {
		return fmt.Errorf("package version is required")
	}

	return nil
}

// GetPythonPackageInfo extracts Python package information
func (p *PythonPostProcessor) GetPythonPackageInfo(cfg config.LLPkgConfig) (string, string, error) {
	if err := p.ValidatePythonPackage(cfg); err != nil {
		return "", "", err
	}

	packageName := cfg.Upstream.Package.Name
	pythonVersion := cfg.Upstream.Package.Version

	// Convert Python version to Go version format
	goVersion := p.convertPythonVersionToGo(pythonVersion)

	return packageName, goVersion, nil
}

// convertPythonVersionToGo converts Python version to Go version format
func (p *PythonPostProcessor) convertPythonVersionToGo(pythonVersion string) string {
	// Ensure version starts with 'v' for Go compatibility
	if !strings.HasPrefix(pythonVersion, "v") {
		return "v" + pythonVersion
	}
	return pythonVersion
}

// CreatePythonRelease creates a release specifically for Python packages
func (p *PythonPostProcessor) CreatePythonRelease(packageName, version string) error {
	releaseTag := fmt.Sprintf("%s/%s", packageName, version)

	// Create release with Python-specific metadata
	release, err := p.DefaultClient.createReleaseByTag(releaseTag)
	if err != nil {
		return fmt.Errorf("failed to create Python release: %v", err)
	}

	// Upload Python-specific artifacts (if any)
	// For Python packages, we might want to upload:
	// - Generated Go bindings
	// - Documentation
	// - Example code
	_, err = p.DefaultClient.uploadArtifactsToRelease(release)
	if err != nil {
		return fmt.Errorf("failed to upload Python artifacts: %v", err)
	}

	return nil
}
