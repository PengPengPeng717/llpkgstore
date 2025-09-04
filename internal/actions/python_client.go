package actions

import (
	"fmt"
	"strings"

	"github.com/goplus/llpkgstore/config"
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
// 直接调用 DefaultClient 的完整实现，确保与 C++ 完全一致
func (p *PythonPostProcessor) Postprocessing() error {
	fmt.Println("Starting Python package post-processing using unified DefaultClient architecture...")

	// 直接调用父类的完整实现，确保所有功能都正常工作
	// 包括：Git 标签创建、GitHub Release 创建、文件上传等
	err := p.DefaultClient.Postprocessing()
	if err != nil {
		return fmt.Errorf("Python package post-processing failed: %v", err)
	}

	fmt.Println("Python package post-processing completed successfully using unified architecture")
	return nil
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
