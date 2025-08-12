package config

import (
	"os"
	"testing"
)

func TestParseLLPkgConfig(t *testing.T) {
	// Test Python package config
	pythonConfig := `{
		"type": "python",
		"upstream": {
			"installer": {
				"name": "pip",
				"config": {
					"python_version": "3.12",
					"index_url": "https://pypi.org/simple/"
				}
			},
			"package": {
				"name": "math",
				"version": "3.12.11"
			}
		}
	}`

	// Write test config to temporary file
	tmpFile, err := os.CreateTemp("", "test-llpkg-*.cfg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(pythonConfig)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Parse config
	cfg, err := ParseLLPkgConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	// Verify Python type
	if cfg.Type != "python" {
		t.Errorf("Expected type 'python', got '%s'", cfg.Type)
	}

	// Verify installer
	if cfg.Upstream.Installer.Name != "pip" {
		t.Errorf("Expected installer 'pip', got '%s'", cfg.Upstream.Installer.Name)
	}

	// Verify package
	if cfg.Upstream.Package.Name != "math" {
		t.Errorf("Expected package 'math', got '%s'", cfg.Upstream.Package.Name)
	}

	if cfg.Upstream.Package.Version != "3.12.11" {
		t.Errorf("Expected version '3.12.11', got '%s'", cfg.Upstream.Package.Version)
	}

	// Test creating upstream
	upstream, err := NewUpstreamFromConfig(cfg.Upstream)
	if err != nil {
		t.Fatalf("Failed to create upstream: %v", err)
	}

	if upstream.Installer.Name() != "pip" {
		t.Errorf("Expected installer name 'pip', got '%s'", upstream.Installer.Name())
	}

	t.Logf("Python package config test passed")
}
