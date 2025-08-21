package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var postProcessingCmd = &cobra.Command{
	Use:   "postprocessing",
	Short: "Process merged PR for Python packages",
	Long:  `Process merged PR for Python packages with simplified post-processing`,
	RunE:  runPythonPostProcessingCmd,
}

func runPythonPostProcessingCmd(_ *cobra.Command, _ []string) error {
	fmt.Println("Starting Python package post-processing...")

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

	// Check for generated Python package files
	requiredFiles := []string{"go.mod", "go.sum"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(currentDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("required file not found: %s", file)
		}
	}

	// For Python packages, use simplified post-processing
	// Directly use v0.0.1 as the version
	version := "v0.0.1"

	fmt.Printf("Python package post-processing completed successfully\n")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Note: This is a simplified post-processing process for Python packages")

	return nil
}

func init() {
	rootCmd.AddCommand(postProcessingCmd)
}
