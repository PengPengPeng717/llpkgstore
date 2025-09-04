package main

// import (
// 	"fmt"
// 	"os"

// 	internal "github.com/goplus/llpkgstore/cmd/llpkgstore/internal_python"
// )

// func main() {
// 	fmt.Println("=== 测试统一命令接口 ===")

// 	// 测试 1: Python 包类型检测
// 	fmt.Println("\n1. 测试 Python 包类型检测...")
// 	testPythonPackageDetection()

// 	// 测试 2: C++ 包类型检测
// 	fmt.Println("\n2. 测试 C++ 包类型检测...")
// 	testCppPackageDetection()

// 	// 测试 3: 统一测试框架
// 	fmt.Println("\n3. 测试统一测试框架...")
// 	testUnifiedTesting()

// 	fmt.Println("\n=== 测试完成 ===")
// }

// func testPythonPackageDetection() {
// 	// 切换到 Python 测试目录
// 	originalDir, _ := os.Getwd()
// 	defer os.Chdir(originalDir)

// 	testDir := "test_python_package"
// 	if err := os.Chdir(testDir); err != nil {
// 		fmt.Printf("❌ 无法切换到测试目录: %v\n", err)
// 		return
// 	}

// 	// 测试包类型检测
// 	processor, err := internal.NewAutoPostProcessor()
// 	if err != nil {
// 		fmt.Printf("❌ 创建自动处理器失败: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("✅ 成功检测到包类型: %s\n", processor.PackageType)
// 	fmt.Printf("✅ 包名称: %s\n", processor.Config.Upstream.Package.Name)
// 	fmt.Printf("✅ 包版本: %s\n", processor.Config.Upstream.Package.Version)
// }

// func testCppPackageDetection() {
// 	// 切换到 C++ 测试目录
// 	originalDir, _ := os.Getwd()
// 	defer os.Chdir(originalDir)

// 	testDir := "test_cpp_package"
// 	if err := os.Chdir(testDir); err != nil {
// 		fmt.Printf("❌ 无法切换到测试目录: %v\n", err)
// 		return
// 	}

// 	// 测试包类型检测
// 	processor, err := internal.NewAutoPostProcessor()
// 	if err != nil {
// 		fmt.Printf("❌ 创建自动处理器失败: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("✅ 成功检测到包类型: %s\n", processor.PackageType)
// 	fmt.Printf("✅ 包名称: %s\n", processor.Config.Upstream.Package.Name)
// 	fmt.Printf("✅ 包版本: %s\n", processor.Config.Upstream.Package.Version)
// }

// func testUnifiedTesting() {
// 	// 切换到 Python 测试目录进行测试
// 	originalDir, _ := os.Getwd()
// 	defer os.Chdir(originalDir)

// 	testDir := "test_python_package"
// 	if err := os.Chdir(testDir); err != nil {
// 		fmt.Printf("❌ 无法切换到测试目录: %v\n", err)
// 		return
// 	}

// 	// 创建一些测试文件
// 	createTestFiles()

// 	// 测试统一测试框架
// 	tester, err := internal.NewUnifiedTester()
// 	if err != nil {
// 		fmt.Printf("❌ 创建统一测试器失败: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("✅ 成功创建统一测试器\n")
// 	fmt.Printf("✅ 包类型: %s\n", tester.PackageType)
// 	fmt.Printf("✅ 包名称: %s\n", tester.PackageName)

// 	// 测试配置验证
// 	if err := tester.Validate(); err != nil {
// 		fmt.Printf("❌ 配置验证失败: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("✅ 配置验证通过\n")
// }

// func createTestFiles() {
// 	// 创建 go.mod 文件
// 	goModContent := `module github.com/test/requests

// go 1.24

// require github.com/goplus/lib v0.2.0
// `
// 	if err := os.WriteFile("go.mod", []byte(goModContent), 0644); err != nil {
// 		fmt.Printf("❌ 创建 go.mod 失败: %v\n", err)
// 		return
// 	}

// 	// 创建 go.sum 文件
// 	goSumContent := `github.com/goplus/lib v0.2.0 h1:example
// `
// 	if err := os.WriteFile("go.sum", []byte(goSumContent), 0644); err != nil {
// 		fmt.Printf("❌ 创建 go.sum 失败: %v\n", err)
// 		return
// 	}

// 	// 创建 _demo 目录和测试文件
// 	if err := os.MkdirAll("_demo/basic_test", 0755); err != nil {
// 		fmt.Printf("❌ 创建 _demo 目录失败: %v\n", err)
// 		return
// 	}

// 	demoContent := `package main

// import (
// 	"fmt"
// 	"github.com/goplus/lib/py"
// 	"requests"
// )

// func main() {
// 	fmt.Println("Testing requests package...")
// 	// Test code would go here
// }
// `
// 	if err := os.WriteFile("_demo/basic_test/main.go", []byte(demoContent), 0644); err != nil {
// 		fmt.Printf("❌ 创建演示文件失败: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("✅ 创建测试文件成功\n")
// }
