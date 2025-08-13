package llpyg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PengPengPeng717/llpkgstore/internal/actions/generator"
	"github.com/PengPengPeng717/llpkgstore/internal/file"
	"github.com/PengPengPeng717/llpkgstore/internal/hashutils"
)

var (
	ErrLLPygGenerate = errors.New("llpyg: cannot generate: ")
	ErrLLPygCheck    = errors.New("llpyg: check fail: ")
)

const (
	// default llpkg repo
	goplusRepo = "github.com/goplus/llpkg/"
	// llpyg running default version
	llpygGoVersion = "1.20.14"
	// llpyg default config file, which MUST exist in specified dir
	llpygConfigFile = "llpyg.cfg"
)

// canHash check file is hashable.
// Hashable file: *.go / llpyg.pub / *.symb.json
func canHash(fileName string) bool {
	if strings.HasSuffix(fileName, ".go") {
		return true
	}
	// Python packages don't have .pub files like C/C++ packages
	// but we can hash other important files
	hashableFiles := map[string]struct{}{
		"go.mod": {},
		"go.sum": {},
	}
	_, ok := hashableFiles[fileName]
	return ok
}

// lockGoVersion locks current Go version to `llpygGoVersion` via GOTOOLCHAIN
func lockGoVersion(cmd *exec.Cmd, pyPath string) {
	// Set Python path for llpyg - use system Python path if pyPath is empty
	if pyPath == "" {
		// Try to find Python site-packages directory
		output, err := exec.Command("python3", "-c", "import site; print(site.getsitepackages()[0])").Output()
		if err == nil {
			pyPath = strings.TrimSpace(string(output))
		} else {
			// Fallback to common paths
			pyPath = "/usr/local/lib/python3.12/site-packages:/usr/lib/python3.12/site-packages"
		}
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("PYTHONPATH=%s", pyPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOTOOLCHAIN=go%s", llpygGoVersion))
}

// diffTwoFiles returns the diff result between a file and b file.
func diffTwoFiles(a, b string) string {
	ret, _ := exec.Command("git", "diff", "--no-index", a, b).CombinedOutput()
	return string(ret)
}

func isExitedUnexpectedly(err error) bool {
	process, ok := err.(*exec.ExitError)
	return ok && !process.Success()
}

// llpygGenerator implements Generator interface, which use llpyg tool to generate llpkg.
type llpygGenerator struct {
	dir         string // llpyg.cfg abs path
	pyDir       string
	packageName string
}

func New(dir, packageName, pyDir string) generator.Generator {
	return &llpygGenerator{dir: dir, packageName: packageName, pyDir: pyDir}
}

// normalizeModulePath returns a normalized module path like
// For llpyg, we need the Python module name, not the Go module path
func (l *llpygGenerator) normalizeModulePath() string {
	return l.packageName
}

func (l *llpygGenerator) findSymbJSON() string {
	matches, _ := filepath.Glob(filepath.Join(l.dir, "*.symb.json"))
	if len(matches) > 0 {
		return filepath.Base(matches[0])
	}
	return ""
}

func (l *llpygGenerator) copyConfigFileTo(path string) error {
	if l.dir == path {
		return nil
	}
	err := file.CopyFile(
		filepath.Join(l.dir, "llpyg.cfg"),
		filepath.Join(path, "llpyg.cfg"),
	)
	// must stop if llpyg.cfg doesn't exist for safety
	if err != nil {
		return err
	}
	if symb := l.findSymbJSON(); symb != "" {
		file.CopyFile(
			filepath.Join(l.dir, symb),
			filepath.Join(path, symb),
		)
	}
	// ignore copy if file doesn't exist
	file.CopyFile(
		filepath.Join(l.dir, "llpyg.pub"),
		filepath.Join(path, "llpyg.pub"),
	)
	return nil
}

func (l *llpygGenerator) Generate(toDir string) error {
	path, err := filepath.Abs(toDir)
	if err != nil {
		return errors.Join(ErrLLPygGenerate, err)
	}
	if err := l.copyConfigFileTo(path); err != nil {
		return errors.Join(ErrLLPygGenerate, err)
	}

	// Create output file for llpyg
	outputFile := filepath.Join(path, l.packageName+".go")
	file, err := os.Create(outputFile)
	if err != nil {
		return errors.Join(ErrLLPygGenerate, err)
	}
	defer file.Close()

	// Execute llpyg command and redirect output to file
	cmd := exec.Command("llpyg", l.normalizeModulePath())
	cmd.Dir = path
	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	// Set environment variables for llpyg
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PYTHONPATH=%s", "/usr/local/lib/python3.12/site-packages:/usr/lib/python3.12/site-packages"),
		fmt.Sprintf("GOTOOLCHAIN=go1.24.5"),
	)

	// llpyg may exit with an error, which may be caused by Stderr.
	// To avoid that case, we have to check its exit code.
	if err := cmd.Run(); isExitedUnexpectedly(err) {
		return errors.Join(ErrLLPygGenerate, err)
	}

	// Generate autogen_link.go file
	if err := l.generateAutogenLinkFile(path); err != nil {
		return errors.Join(ErrLLPygGenerate, err)
	}

	// Generate go.mod file if it doesn't exist
	if err := l.generateGoModFile(path); err != nil {
		return errors.Join(ErrLLPygGenerate, err)
	}

	return nil
}

func (l *llpygGenerator) Check(dir string) error {
	baseDir, err := filepath.Abs(dir)
	if err != nil {
		return errors.Join(ErrLLPygCheck, err)
	}

	// 1. compute hash
	generated, err := hashutils.Dir(baseDir, canHash)
	if err != nil {
		return errors.Join(ErrLLPygCheck, err)
	}
	userGenerated, err := hashutils.Dir(l.dir, canHash)
	if err != nil {
		return errors.Join(ErrLLPygCheck, err)
	}

	// 2. check hash
	for name, hash := range userGenerated {
		generatedHash, ok := generated[name]
		if !ok {
			// if this file is hashable, it's unexpected
			// if not, we can skip it safely.
			if canHash(name) {
				return errors.Join(ErrLLPygCheck, fmt.Errorf("unexpected file: %s", name))
			}
			// skip file
			continue
		}
		if !strings.EqualFold(string(hash), string(generatedHash)) {
			return errors.Join(ErrLLPygCheck, fmt.Errorf("file not equal: %s %s", name,
				diffTwoFiles(filepath.Join(l.dir, name), filepath.Join(baseDir, name))))
		}
	}
	// 3. check missing file
	for name := range generated {
		if _, ok := userGenerated[name]; !ok {
			return errors.Join(ErrLLPygCheck, fmt.Errorf("missing file: %s", name))
		}
	}
	return nil
}

func (l *llpygGenerator) generateAutogenLinkFile(path string) error {
	content := fmt.Sprintf(`package %s

import _ "github.com/goplus/lib/py"

// LLGoPackage is defined in the main binding file
`, l.packageName)

	autogenFile := filepath.Join(path, l.packageName+"_autogen_link.go")
	return os.WriteFile(autogenFile, []byte(content), 0644)
}

func (l *llpygGenerator) generateGoModFile(path string) error {
	goModFile := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModFile); err == nil {
		// go.mod already exists, skip
		return nil
	}

	content := fmt.Sprintf(`module github.com/goplus/llpkg/%s

go 1.23

require github.com/goplus/lib v0.2.0
`, l.packageName)

	return os.WriteFile(goModFile, []byte(content), 0644)
}
