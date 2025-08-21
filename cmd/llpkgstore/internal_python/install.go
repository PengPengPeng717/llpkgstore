package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/goplus/llpkgstore/config"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [LLPkgConfigFilePath]",
	Short: "Manually install a Python package",
	Long:  `Manually install a Python package from llpkg.cfg file using pip.`,
	Args:  cobra.ExactArgs(1),
	RunE:  manuallyInstall,
}

func manuallyInstall(cmd *cobra.Command, args []string) error {
	cfgPath := args[0]

	// 检查配置文件是否存在
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", cfgPath)
	}

	// 解析配置文件
	LLPkgConfig, err := config.ParseLLPkgConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 检查包类型
	if LLPkgConfig.Type != "python" {
		return fmt.Errorf("不支持的包类型: %s，当前仅支持 Python 包", LLPkgConfig.Type)
	}

	// 获取输出目录
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	// 如果输出目录为空，使用当前目录
	if output == "" {
		output = "."
	}

	// 确保输出目录存在
	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	log.Printf("开始安装 Python 包: %s==%s", LLPkgConfig.Upstream.Package.Name, LLPkgConfig.Upstream.Package.Version)
	log.Printf("输出目录: %s", output)

	// 创建 upstream 实例
	upstream, err := config.NewUpstreamFromConfig(LLPkgConfig.Upstream)
	if err != nil {
		return fmt.Errorf("创建 upstream 实例失败: %v", err)
	}

	// 执行安装
	installedPackages, err := upstream.Installer.Install(upstream.Pkg, output)
	if err != nil {
		return fmt.Errorf("安装失败: %v", err)
	}

	log.Printf("安装成功！已安装的包: %v", installedPackages)

	// 显示安装结果
	fmt.Printf("✓ 成功安装 Python 包: %s==%s\n", LLPkgConfig.Upstream.Package.Name, LLPkgConfig.Upstream.Package.Version)
	fmt.Printf("  安装位置: %s\n", output)
	fmt.Printf("  已安装的包: %v\n", installedPackages)

	// 如果是 pip 安装器，显示额外的信息
	if LLPkgConfig.Upstream.Installer.Name == "pip" {
		fmt.Println("\n注意:")
		fmt.Println("- 包已通过 pip3 安装到指定目录")
		fmt.Println("- 可以使用 'generate' 命令生成 Go 绑定")
		fmt.Println("- 可以使用 'test' 命令验证安装结果")
	}

	return nil
}

func init() {
	installCmd.Flags().StringP("output", "o", "", "安装输出目录 (默认: 当前目录)")
	rootCmd.AddCommand(installCmd)
}
