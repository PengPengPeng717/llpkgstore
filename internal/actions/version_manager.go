package actions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goplus/llpkgstore/internal/actions/env"
	"github.com/goplus/llpkgstore/internal/actions/versions"
)

// VersionManager provides unified version management for all package types
type VersionManager struct {
	*DefaultClient
}

// NewVersionManager creates a new version manager
func NewVersionManager() (*VersionManager, error) {
	client, err := NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create default client: %v", err)
	}

	return &VersionManager{
		DefaultClient: client,
	}, nil
}

// ExtractVersionFromCommit extracts version from commit message using unified logic
func (vm *VersionManager) ExtractVersionFromCommit() (string, error) {
	sha, err := env.LatestCommitSHA()
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit SHA: %v", err)
	}

	commit, err := vm.commitMessage(sha)
	if err != nil {
		return "", fmt.Errorf("failed to get commit message: %v", err)
	}

	message := commit.GetCommit().GetMessage()
	return vm.parseVersionFromMessage(message)
}

// parseVersionFromMessage parses version from commit message using unified patterns
func (vm *VersionManager) parseVersionFromMessage(message string) (string, error) {
	// Unified patterns for all package types
	patterns := []string{
		`Release-as:\s*[^/]+/(v[\d.]+)`, // "Release-as: numpy/v1.26.4"
		`Release-as:\s*(v[\d.]+)`,       // "Release-as: v1.26.4"
		`Release:\s*[^/]+/(v[\d.]+)`,    // "Release: numpy/v1.26.4"
		`Release:\s*(v[\d.]+)`,          // "Release: v1.26.4"
		`Version:\s*(v[\d.]+)`,          // "Version: v1.26.4"
		`version:\s*(v[\d.]+)`,          // "version: v1.26.4"
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 1 {
			version := matches[1]
			if err := vm.validateVersion(version); err != nil {
				continue
			}
			return version, nil
		}
	}

	return "", fmt.Errorf("no version pattern found in commit message")
}

// validateVersion validates version format
func (vm *VersionManager) validateVersion(version string) error {
	if !strings.HasPrefix(version, "v") {
		return fmt.Errorf("version must start with 'v': %s", version)
	}

	// Basic semver validation
	if !vm.isValidSemver(version) {
		return fmt.Errorf("invalid semver format: %s", version)
	}

	return nil
}

// isValidSemver checks if version follows semver format
func (vm *VersionManager) isValidSemver(version string) bool {
	// Remove 'v' prefix for validation
	version = strings.TrimPrefix(version, "v")

	// Basic semver pattern: major.minor.patch[-prerelease][+build]
	pattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

// UpdateVersionMapping updates version mapping for any package type
func (vm *VersionManager) UpdateVersionMapping(packageName, upstreamVersion, goVersion string) error {
	ver := versions.Read("llpkgstore.json")
	ver.Write(packageName, upstreamVersion, goVersion)
	return nil
}

// GetVersionMapping retrieves version mapping for a package
func (vm *VersionManager) GetVersionMapping(packageName string) ([]string, error) {
	ver := versions.Read("llpkgstore.json")
	return ver.CVersions(packageName), nil
}

// ConvertVersion converts upstream version to Go version format
func (vm *VersionManager) ConvertVersion(upstreamVersion string) string {
	// Ensure version starts with 'v' for Go compatibility
	if !strings.HasPrefix(upstreamVersion, "v") {
		return "v" + upstreamVersion
	}
	return upstreamVersion
}

// ParseMappedVersion parses mapped version string
func (vm *VersionManager) ParseMappedVersion(version string) (string, string, error) {
	// Format: "packageName/v1.0.0"
	parts := strings.Split(version, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid mapped version format: %s", version)
	}

	packageName := parts[0]
	mappedVersion := parts[1]

	return packageName, mappedVersion, nil
}

// CreateVersionTag creates a version tag for any package type
func (vm *VersionManager) CreateVersionTag(packageName, version string) (string, error) {
	// Format: "packageName/version"
	tag := fmt.Sprintf("%s/%s", packageName, version)

	// Validate tag format
	if err := vm.validateTag(tag); err != nil {
		return "", fmt.Errorf("invalid tag format: %v", err)
	}

	return tag, nil
}

// validateTag validates tag format
func (vm *VersionManager) validateTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}

	if strings.Contains(tag, " ") {
		return fmt.Errorf("tag cannot contain spaces: %s", tag)
	}

	return nil
}

// HasVersionTag checks if a version tag already exists
func (vm *VersionManager) HasVersionTag(tag string) bool {
	return hasTag(tag)
}

// CreateVersionRelease creates a release for any package type
func (vm *VersionManager) CreateVersionRelease(packageName, version string) error {
	tag, err := vm.CreateVersionTag(packageName, version)
	if err != nil {
		return fmt.Errorf("failed to create version tag: %v", err)
	}

	// Check if tag already exists
	if vm.HasVersionTag(tag) {
		return fmt.Errorf("tag already exists: %s", tag)
	}

	// Get commit SHA
	sha, err := env.LatestCommitSHA()
	if err != nil {
		return fmt.Errorf("failed to get commit SHA: %v", err)
	}

	// Create tag
	if err := vm.DefaultClient.createTag(tag, sha); err != nil {
		return fmt.Errorf("failed to create tag: %v", err)
	}

	// Create release
	release, err := vm.DefaultClient.createReleaseByTag(tag)
	if err != nil {
		return fmt.Errorf("failed to create release: %v", err)
	}

	// Upload artifacts
	_, err = vm.DefaultClient.uploadArtifactsToRelease(release)
	if err != nil {
		return fmt.Errorf("failed to upload artifacts: %v", err)
	}

	return nil
}

// CleanupLegacyBranch handles legacy branch cleanup for any package type
func (vm *VersionManager) CleanupLegacyBranch() error {
	branchName, isLegacy, err := vm.DefaultClient.isLegacyVersion()
	if err != nil {
		return err
	}

	if isLegacy {
		fmt.Printf("Cleaning up legacy branch: %s\n", branchName)
		return vm.DefaultClient.removeBranch(branchName)
	}

	return nil
}
