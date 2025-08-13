package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/PengPengPeng717/llpkgstore/config"
	"github.com/PengPengPeng717/llpkgstore/internal/actions/generator"
	"github.com/PengPengPeng717/llpkgstore/internal/actions/generator/llcppg"
	"github.com/PengPengPeng717/llpkgstore/internal/actions/generator/llpyg"
	"github.com/PengPengPeng717/llpkgstore/internal/file"
	"github.com/PengPengPeng717/llpkgstore/internal/pc"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate LLGo bindings",
	Long:  ``,
	RunE:  runGenerate,
}

func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func runGenerateWithDir(dir string) error {
	cfg, err := config.ParseLLPkgConfig(filepath.Join(dir, LLGOModuleIdentifyFile))
	if err != nil {
		return fmt.Errorf("parse config error: %v", err)
	}
	uc, err := config.NewUpstreamFromConfig(cfg.Upstream)
	if err != nil {
		return err
	}
	log.Printf("Start to generate %s", uc.Pkg.Name)

	tempDir, err := os.MkdirTemp("", "llpkg-tool")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	var pcName []string
	// Skip installation for Python builtin modules
	if cfg.Type == "python" && uc.Pkg.Version == "builtin" {
		pcName = []string{uc.Pkg.Name}
	} else {
		pcName, err = uc.Installer.Install(uc.Pkg, tempDir)
		if err != nil {
			return err
		}
	}

	// copy file for debugging (only for C/C++ packages)
	if cfg.Type != "python" {
		err = file.CopyFilePattern(tempDir, dir, "*.pc")
		if err != nil {
			return err
		}
	}

	var gen generator.Generator

	// Check if it's a Python package
	if cfg.Type == "python" {
		// try llpygcfg if llpyg.cfg doesn't exist
		if _, err := os.Stat(filepath.Join(dir, "llpyg.cfg")); os.IsNotExist(err) {
			cmd := exec.Command("llpygcfg", uc.Pkg.Name)
			cmd.Dir = dir
			ret, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("llpygcfg execute fail: %s", string(ret))
			}
		}
		gen = llpyg.New(dir, cfg.Upstream.Package.Name, tempDir)
	} else {
		// try llcppcfg if llcppg.cfg doesn't exist
		if _, err := os.Stat(filepath.Join(dir, "llcppg.cfg")); os.IsNotExist(err) {
			cmd := exec.Command("llcppcfg", pcName[0])
			cmd.Dir = dir
			pc.SetPath(cmd, tempDir)
			ret, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("llcppcfg execute fail: %s", string(ret))
			}
		}
		gen = llcppg.New(dir, cfg.Upstream.Package.Name, tempDir)
	}

	return gen.Generate(dir)
}

func runGenerate(_ *cobra.Command, args []string) error {
	exec.Command("conan", "profile", "detect").Run()

	path := currentDir()
	// by default, use current dir
	if len(args) == 0 {
		return runGenerateWithDir(path)
	}
	for _, argPath := range args {
		absPath, err := filepath.Abs(argPath)
		if err != nil {
			continue
		}
		err = runGenerateWithDir(absPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
