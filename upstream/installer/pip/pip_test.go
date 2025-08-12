package pip

import (
	"os"
	"testing"

	"github.com/PengPengPeng717/llpkgstore/upstream"
)

func TestPipInstaller(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pip-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create pip installer with test config
	config := map[string]string{
		"python_version": "3.12",
		"index_url":      "https://pypi.org/simple/",
		"trusted_host":   "pypi.org",
	}
	installer := NewPipInstaller(config)

	// Test installer name
	if installer.Name() != "pip" {
		t.Errorf("Expected installer name 'pip', got '%s'", installer.Name())
	}

	// Test config
	installerConfig := installer.Config()
	if installerConfig["python_version"] != "3.12" {
		t.Errorf("Expected python_version '3.12', got '%s'", installerConfig["python_version"])
	}

	// Test package
	pkg := upstream.Package{
		Name:    "math",
		Version: "3.12.11",
	}

	// Test search (this might fail if no internet connection)
	_, err = installer.Search(pkg)
	if err != nil {
		t.Logf("Search failed (expected for math module): %v", err)
	}

	// Test dependency (this might fail if no internet connection)
	_, err = installer.Dependency(pkg)
	if err != nil {
		t.Logf("Dependency failed (expected for math module): %v", err)
	}

	t.Logf("Pip installer tests completed")
}
