package pip

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PengPengPeng717/llpkgstore/upstream"
)

var (
	ErrPackageNotFound = errors.New("python package not found")
	ErrModuleNotFound  = errors.New("python module not found")
)

type pipInstaller struct {
	config map[string]string
}

// NewPipInstaller creates a new pip installer with the given configuration
func NewPipInstaller(config map[string]string) upstream.Installer {
	return &pipInstaller{config: config}
}

func (p *pipInstaller) Name() string {
	return "pip"
}

func (p *pipInstaller) Config() map[string]string {
	return p.config
}

// Install downloads and installs the specified Python package
func (p *pipInstaller) Install(pkg upstream.Package, outputDir string) ([]string, error) {
	// Create requirements.txt for the package
	requirementsFile := filepath.Join(outputDir, "requirements.txt")
	requirementsContent := fmt.Sprintf("%s==%s", pkg.Name, pkg.Version)

	err := os.WriteFile(requirementsFile, []byte(requirementsContent), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create requirements.txt: %v", err)
	}

	// Build pip install command
	cmd := exec.Command("pip3", "install", "-r", requirementsFile, "--target", outputDir)

	// Add pip configuration options
	if indexURL := p.config["index_url"]; indexURL != "" {
		cmd.Args = append(cmd.Args, "-i", indexURL)
	}
	if extraIndexURL := p.config["extra_index_url"]; extraIndexURL != "" {
		cmd.Args = append(cmd.Args, "--extra-index-url", extraIndexURL)
	}
	if trustedHost := p.config["trusted_host"]; trustedHost != "" {
		cmd.Args = append(cmd.Args, "--trusted-host", trustedHost)
	}

	// Set Python version if specified
	if pythonVersion := p.config["python_version"]; pythonVersion != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PYTHON_VERSION=%s", pythonVersion))
	}

	// Execute pip install
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("pip install failed: %v, output: %s", err, string(output))
	}

	// For Python packages, we return the module name as the "pkg-config" equivalent
	// This allows the generator to locate the installed module
	return []string{pkg.Name}, nil
}

// Search checks PyPI for the specified package availability
func (p *pipInstaller) Search(pkg upstream.Package) ([]string, error) {
	cmd := exec.Command("pip3", "search", pkg.Name)

	// Add pip configuration options
	if indexURL := p.config["index_url"]; indexURL != "" {
		cmd.Args = append(cmd.Args, "-i", indexURL)
	}
	if trustedHost := p.config["trusted_host"]; trustedHost != "" {
		cmd.Args = append(cmd.Args, "--trusted-host", trustedHost)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("pip search failed: %v", err)
	}

	// Parse search results
	lines := strings.Split(string(output), "\n")
	var results []string
	for _, line := range lines {
		if strings.Contains(line, pkg.Name) {
			results = append(results, strings.TrimSpace(line))
		}
	}

	return results, nil
}

// Dependency retrieves the list of dependencies for the specified Python package
func (p *pipInstaller) Dependency(pkg upstream.Package) ([]upstream.Package, error) {
	// Create temporary directory for dependency analysis
	tempDir, err := os.MkdirTemp("", "pip-deps-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Install package to temp directory to analyze dependencies
	_, err = p.Install(pkg, tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to install package for dependency analysis: %v", err)
	}

	// Use pip show to get dependency information
	cmd := exec.Command("pip3", "show", pkg.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("pip show failed: %v", err)
	}

	// Parse dependencies from pip show output
	var dependencies []upstream.Package
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Requires:") {
			requires := strings.TrimSpace(strings.TrimPrefix(line, "Requires:"))
			if requires != "" {
				// Parse individual dependencies
				deps := strings.Split(requires, ",")
				for _, dep := range deps {
					dep = strings.TrimSpace(dep)
					if dep != "" {
						// Extract package name and version
						parts := strings.Split(dep, " ")
						if len(parts) >= 1 {
							dependencies = append(dependencies, upstream.Package{
								Name:    parts[0],
								Version: "", // Version info might need additional parsing
							})
						}
					}
				}
			}
			break
		}
	}

	return dependencies, nil
}
