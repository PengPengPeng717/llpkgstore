package config

import (
	"errors"

	"github.com/goplus/llpkgstore/upstream"
	"github.com/goplus/llpkgstore/upstream/installer/conan"
	"github.com/goplus/llpkgstore/upstream/installer/pip"
)

var ValidInstallers = []string{"conan", "pip"}

// LLPkgConfig represents the configuration structure parsed from llpkg.cfg files.
type LLPkgConfig struct {
	Type     string         `json:"type,omitempty"` // "python" for Python packages, empty for C/C++
	Upstream UpstreamConfig `json:"upstream"`
	Llpyg    LlpygConfig    `json:"llpyg,omitempty"` // llpyg command line options
}

// LlpygConfig defines the llpyg command line options
type LlpygConfig struct {
	OutputDir string `json:"output_dir,omitempty"` // -o option, default "./test"
	ModName   string `json:"mod_name,omitempty"`   // -mod option, default package name
	ModDepth  int    `json:"mod_depth,omitempty"`  // -d option, default 1
}

// Validate validates the LlpygConfig
func (l *LlpygConfig) Validate() error {
	if l.ModDepth < 0 {
		return errors.New("mod_depth must be non-negative")
	}
	if l.ModDepth > 10 {
		return errors.New("mod_depth should not exceed 10 for performance reasons")
	}
	return nil
}

// GetDefaultModDepth returns the default module depth if not specified
func (l *LlpygConfig) GetDefaultModDepth() int {
	if l.ModDepth == 0 {
		return 1 // 默认深度为 1
	}
	return l.ModDepth
}

// GetDefaultOutputDir returns the default output directory if not specified
func (l *LlpygConfig) GetDefaultOutputDir() string {
	if l.OutputDir == "" {
		return "./test" // 默认输出目录
	}
	return l.OutputDir
}

// UpstreamConfig defines the upstream configuration containing installer settings and package metadata.
type UpstreamConfig struct {
	Installer InstallerConfig `json:"installer"`
	Package   PackageConfig   `json:"package"`
}

// InstallerConfig specifies the installer type and its configuration options.
// "name" field must match supported installers (e.g., "conan").
// "config" holds installer-specific parameters (optional).
type InstallerConfig struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config,omitempty"`
}

// PackageConfig defines the target library package's identifier and version requirements.
type PackageConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewUpstreamFromConfig creates an Upstream instance from configuration data.
// Returns error if unsupported installer type is specified.
func NewUpstreamFromConfig(upstreamConfig UpstreamConfig) (*upstream.Upstream, error) {
	switch upstreamConfig.Installer.Name {
	case "conan":
		return &upstream.Upstream{
			Installer: conan.NewConanInstaller(upstreamConfig.Installer.Config),
			Pkg: upstream.Package{
				Name:    upstreamConfig.Package.Name,
				Version: upstreamConfig.Package.Version,
			},
		}, nil
	case "pip":
		return &upstream.Upstream{
			Installer: pip.NewPipInstaller(upstreamConfig.Installer.Config),
			Pkg: upstream.Package{
				Name:    upstreamConfig.Package.Name,
				Version: upstreamConfig.Package.Version,
			},
		}, nil
	default:
		return nil, errors.New("unknown upstream installer: " + upstreamConfig.Installer.Name)
	}
}
