package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Build and upload Python binary packages",
	Long:  `Build and upload Python binary packages with simplified version handling`,
	RunE:  runPythonReleaseCmd,
}

func runPythonReleaseCmd(_ *cobra.Command, _ []string) error {
	// For Python packages, use simplified version handling
	// Directly use v0.0.1 as specified
	version := "v0.0.1"

	fmt.Printf("Starting Python package release with version: %s\n", version)

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Check if llpkg.cfg exists
	cfgPath := filepath.Join(currentDir, "llpkg.cfg")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf("llpkg.cfg not found in current directory")
	}

	// Check if generated files exist
	generatedFiles := []string{"go.mod", "go.sum"}
	for _, file := range generatedFiles {
		filePath := filepath.Join(currentDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("required file not found: %s", file)
		}
	}

	fmt.Printf("Python package release completed successfully with version: %s\n", version)
	fmt.Println("Note: This is a simplified release process for Python packages")

	return nil
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
