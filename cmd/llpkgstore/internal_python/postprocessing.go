package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/llpkgstore/config"
	"github.com/spf13/cobra"
)

var postProcessingCmd = &cobra.Command{
	Use:   "postprocessing",
	Short: "Process merged PR for Python packages",
	Long:  `Process merged PR for Python packages with GitHub Release support`,
	RunE:  runPythonPostProcessingCmd,
}

// LLPkgStoreJSON represents the structure of llpkgstore.json
type LLPkgStoreJSON struct {
	Packages map[string]PackageInfo `json:"packages"`
}

// PackageInfo represents package version information
type PackageInfo struct {
	Versions []VersionInfo `json:"versions"`
}

// VersionInfo represents version mapping
type VersionInfo struct {
	Python string   `json:"python"`
	Go     []string `json:"go"`
}

func runPythonPostProcessingCmd(_ *cobra.Command, _ []string) error {
	fmt.Println("Starting Python package post-processing...")

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Check if llpkg.cfg exists and parse it
	cfgPath := filepath.Join(currentDir, "llpkg.cfg")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf("llpkg.cfg not found in current directory")
	}

	cfg, err := config.ParseLLPkgConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to parse llpkg.cfg: %v", err)
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
	pythonVersion := cfg.Upstream.Package.Version
	packageName := cfg.Upstream.Package.Name

	// Generate or update llpkgstore.json
	if err := updateLLPkgStoreJSON(packageName, pythonVersion, version); err != nil {
		return fmt.Errorf("failed to update llpkgstore.json: %v", err)
	}

	// Try to create GitHub Release if we're in a GitHub Actions environment
	if err := createGitHubRelease(packageName, version, currentDir); err != nil {
		fmt.Printf("Warning: Failed to create GitHub Release: %v\n", err)
		fmt.Println("This is normal if not running in GitHub Actions or if release already exists")
	}

	fmt.Printf("Python package post-processing completed successfully\n")
	fmt.Printf("Package: %s\n", packageName)
	fmt.Printf("Python Version: %s\n", pythonVersion)
	fmt.Printf("Go Version: %s\n", version)
	fmt.Println("Note: This is a simplified post-processing process for Python packages")

	return nil
}

// createGitHubRelease attempts to create a GitHub Release for the package
func createGitHubRelease(packageName, version, currentDir string) error {
	// Check if we're in a GitHub Actions environment
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return fmt.Errorf("not running in GitHub Actions environment")
	}

	// Create a simple release for Python packages
	// Since Python packages don't follow the same version mapping as C++ packages,
	// we'll create a release with the package name and version
	releaseTag := fmt.Sprintf("%s-%s", packageName, version)

	// Use GitHub CLI to create the release
	// First check if release already exists
	checkCmd := fmt.Sprintf("gh release view %s --repo $GITHUB_REPOSITORY >/dev/null 2>&1", releaseTag)
	if err := exec.Command("bash", "-c", checkCmd).Run(); err == nil {
		fmt.Printf("Release %s already exists, skipping creation\n", releaseTag)
		return nil
	}

	// Create the release using GitHub CLI
	createCmd := fmt.Sprintf("gh release create %s --title 'Release %s %s' --notes 'Automated release for %s version %s' --repo $GITHUB_REPOSITORY --draft=false --prerelease=false",
		releaseTag, packageName, version, packageName, version)

	cmd := exec.Command("bash", "-c", createCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create GitHub release: %v", err)
	}

	fmt.Printf("Successfully created GitHub Release: %s\n", releaseTag)
	fmt.Printf("Package %s version %s is now available via: llgo get github.com/$GITHUB_REPOSITORY/%s@%s\n",
		packageName, version, packageName, version)

	return nil
}

// updateLLPkgStoreJSON updates the llpkgstore.json file with new package information
func updateLLPkgStoreJSON(packageName, pythonVersion, goVersion string) error {
	jsonPath := "llpkgstore.json"

	// Read existing JSON file if it exists
	var llpkgStore LLPkgStoreJSON
	if _, err := os.Stat(jsonPath); err == nil {
		// File exists, read it
		data, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("failed to read existing llpkgstore.json: %v", err)
		}

		if err := json.Unmarshal(data, &llpkgStore); err != nil {
			return fmt.Errorf("failed to parse existing llpkgstore.json: %v", err)
		}
	} else {
		// File doesn't exist, create new structure
		llpkgStore = LLPkgStoreJSON{
			Packages: make(map[string]PackageInfo),
		}
	}

	// Initialize package info if it doesn't exist
	if _, exists := llpkgStore.Packages[packageName]; !exists {
		llpkgStore.Packages[packageName] = PackageInfo{
			Versions: []VersionInfo{},
		}
	}

	// Check if this Python version already exists
	packageInfo := llpkgStore.Packages[packageName]
	found := false
	for i, versionInfo := range packageInfo.Versions {
		if versionInfo.Python == pythonVersion {
			// Python version exists, add Go version if not already present
			if !contains(versionInfo.Go, goVersion) {
				packageInfo.Versions[i].Go = append(versionInfo.Go, goVersion)
			}
			found = true
			break
		}
	}

	// If Python version doesn't exist, add new version entry
	if !found {
		packageInfo.Versions = append(packageInfo.Versions, VersionInfo{
			Python: pythonVersion,
			Go:     []string{goVersion},
		})
	}

	llpkgStore.Packages[packageName] = packageInfo

	// Write updated JSON back to file
	data, err := json.MarshalIndent(llpkgStore, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal llpkgstore.json: %v", err)
	}

	if err := os.WriteFile(jsonPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write llpkgstore.json: %v", err)
	}

	fmt.Printf("Updated llpkgstore.json with package %s (Python: %s -> Go: %s)\n", packageName, pythonVersion, goVersion)
	return nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(postProcessingCmd)
}
