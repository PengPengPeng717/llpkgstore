package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goplus/llpkgstore/config"
	"github.com/goplus/llpkgstore/internal/actions"
	"github.com/spf13/cobra"
)

// 重新定义 postprocessing 命令，使用与 C++ 一致的架构
var postProcessingCmd = &cobra.Command{
	Use:   "postprocessing",
	Short: "Process merged PR with unified architecture",
	Long:  `Process merged PR using unified DefaultClient architecture, automatically detects package type and applies consistent version management`,
	RunE:  runUnifiedPostProcessingCmd,
}

// UnifiedPostProcessor 使用与 C++ 一致的架构处理所有包类型
type UnifiedPostProcessor struct {
	PackageType string
	Config      config.LLPkgConfig
}

// NewUnifiedPostProcessor 创建统一的后处理器
func NewUnifiedPostProcessor() (*UnifiedPostProcessor, error) {
	// 从 llpkg.cfg 检测包类型
	cfg, packageType, err := detectPackageType(".")
	if err != nil {
		return nil, fmt.Errorf("failed to detect package type: %v", err)
	}

	return &UnifiedPostProcessor{
		PackageType: packageType,
		Config:      cfg,
	}, nil
}

// Postprocessing 使用统一架构处理后处理逻辑
func (u *UnifiedPostProcessor) Postprocessing() error {
	fmt.Printf("Detected package type: %s\n", u.PackageType)
	fmt.Printf("Package: %s\n", u.Config.Upstream.Package.Name)
	fmt.Printf("Version: %s\n", u.Config.Upstream.Package.Version)

	switch u.PackageType {
	case "python":
		return u.processPythonPackage()
	case "cpp":
		return u.processCppPackage()
	default:
		return fmt.Errorf("unsupported package type: %s", u.PackageType)
	}
}

// processPythonPackage 使用统一的 DefaultClient 接口处理 Python 包
func (u *UnifiedPostProcessor) processPythonPackage() error {
	fmt.Println("Processing Python package using unified DefaultClient architecture...")

	// 创建 Python 后处理器，使用与 C++ 一致的接口
	processor, err := actions.NewPythonPostProcessor()
	if err != nil {
		return fmt.Errorf("failed to create Python processor: %v", err)
	}

	// 使用统一的后处理逻辑
	return processor.Postprocessing()
}

// processCppPackage 使用现有的 DefaultClient 接口处理 C++ 包
func (u *UnifiedPostProcessor) processCppPackage() error {
	fmt.Println("Processing C++ package using existing DefaultClient interface...")

	// 创建默认客户端用于 C++ 包
	client, err := actions.NewDefaultClient()
	if err != nil {
		return fmt.Errorf("failed to create default client: %v", err)
	}

	// 使用现有的后处理逻辑
	return client.Postprocessing()
}

// detectPackageType 从 llpkg.cfg 检测包类型
func detectPackageType(dir string) (config.LLPkgConfig, string, error) {
	cfgPath := filepath.Join(dir, "llpkg.cfg")

	// 检查 llpkg.cfg 是否存在
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return config.LLPkgConfig{}, "", fmt.Errorf("llpkg.cfg not found in directory: %s", dir)
	}

	// 解析配置
	cfg, err := config.ParseLLPkgConfig(cfgPath)
	if err != nil {
		return config.LLPkgConfig{}, "", fmt.Errorf("failed to parse llpkg.cfg: %v", err)
	}

	// 确定包类型
	packageType := cfg.Type
	if packageType == "" {
		// 如果类型未指定，默认为 cpp
		packageType = "cpp"
	}

	return cfg, packageType, nil
}

func runUnifiedPostProcessingCmd(_ *cobra.Command, _ []string) error {
	processor, err := NewUnifiedPostProcessor()
	if err != nil {
		return err
	}
	return processor.Postprocessing()
}

func init() {
	rootCmd.AddCommand(postProcessingCmd)
}
