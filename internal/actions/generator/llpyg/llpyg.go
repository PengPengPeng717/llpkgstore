package llpyg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/llpkgstore/internal/actions/generator"
	"github.com/goplus/llpkgstore/internal/file"
	"github.com/goplus/llpkgstore/internal/hashutils"
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
	// Set Python path for llpyg
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
// math => github.com/goplus/llpkg/math
func (l *llpygGenerator) normalizeModulePath() string {
	return goplusRepo + l.packageName
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

	// Execute llpyg command
	cmd := exec.Command("llpyg", "-mod", l.normalizeModulePath(), llpygConfigFile)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	lockGoVersion(cmd, l.pyDir)

	// llpyg may exit with an error, which may be caused by Stderr.
	// To avoid that case, we have to check its exit code.
	if err := cmd.Run(); isExitedUnexpectedly(err) {
		return errors.Join(ErrLLPygGenerate, err)
	}

	// check output again
	generatedPath := filepath.Join(path, l.packageName)
	if _, err := os.Stat(generatedPath); os.IsNotExist(err) {
		return errors.Join(ErrLLPygCheck, errors.New("generate fail"))
	}

	// copy out the generated result
	// be careful: llpyg result MUST not override existed file,
	// otherwise, checking is meaningless.
	err = file.CopyFS(path, os.DirFS(generatedPath), true)
	if err != nil {
		return errors.Join(ErrLLPygGenerate, err)
	}

	os.RemoveAll(generatedPath)
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
