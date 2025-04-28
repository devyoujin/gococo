package coverage

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	moduleCoverageDir = "modules"
)

type module struct {
	path         string
	name         string
	coverageData []byte
}

func (module *module) generateCoverage() error {
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d.out", time.Now().UnixNano()))
	cmd := exec.Command("go", "test", "./...", "-coverprofile="+tempFile)
	cmd.Dir = module.path
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run tests in %s: %w", module.path, err)
	}
	content, err := os.ReadFile(tempFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", tempFile, err)
	}
	module.coverageData = content
	return nil
}

func (module *module) saveCoverageToFile(outputDir string) error {
	if module.coverageData == nil {
		return fmt.Errorf("no coverage data available in %s", module.name)
	}
	outputPath := filepath.Join(outputDir, moduleCoverageDir, module.name+".coverage.out")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	if err := os.WriteFile(outputPath, module.coverageData, 0644); err != nil {
		return fmt.Errorf("failed to write coverage data to file in %s", module.name)
	}
	return nil
}

func findGoModules(root string) ([]module, error) {
	var modules []module
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || entry.Name() != "go.mod" {
			return nil
		}
		modPath := filepath.Dir(path)
		modName := filepath.Base(modPath)
		modules = append(modules, module{name: modName, path: modPath})
		return nil
	})
	return modules, err
}
