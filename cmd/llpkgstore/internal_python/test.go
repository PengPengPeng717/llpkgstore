package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Python verification functionality",
	Long:  `Test the verification and check functionality for Python packages`,
	RunE:  runTest,
}

func runTest(_ *cobra.Command, args []string) error {
	// 确定测试目录
	testDir := "."
	if len(args) > 0 {
		if absPath, err := filepath.Abs(args[0]); err == nil {
			testDir = absPath
		}
	}

	// 检查测试目录是否存在
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		return fmt.Errorf("测试目录不存在: %s", testDir)
	}

	// 检查 llpkg.cfg 文件是否存在
	cfgPath := filepath.Join(testDir, LLGOModuleIdentifyFile)
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", cfgPath)
	}

	fmt.Printf("开始测试目录: %s\n", testDir)

	// 运行验证测试
	err := TestVerification(testDir)
	if err != nil {
		return fmt.Errorf("验证测试失败: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(testCmd)
}
