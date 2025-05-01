package coverage

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
)

type ManagerInterface interface {
	FindGoModules() ([]module, error)
	GenerateCoverages(modules []module) error
	GenerateCoverProfile() error
}

type Manager struct {
	pwd string
	mergedCoverageDir string
	coverageProfile string
}

func NewManager(
	pwd string,
	mergedCoverageDir string,
	coverageProfile string) ManagerInterface {
	return &Manager{
		pwd:               pwd,
		mergedCoverageDir: mergedCoverageDir,
		coverageProfile:  coverageProfile,
	}
}

func (manager *Manager) FindGoModules() ([]module, error) {
	var modules []module
	err := filepath.WalkDir(manager.pwd, func(path string, entry fs.DirEntry, err error) error {
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

func (manager *Manager) GenerateCoverages(modules []module) error {
	for _, module := range modules {
		if err := module.generateCoverage(manager.mergedCoverageDir); err != nil {
			return fmt.Errorf("failed to generate coverage in %s: %w", module.name, err)
		}
	}
	return nil
}

func (manager *Manager) GenerateCoverProfile() error {
	cmd := exec.Command("go", "tool", "covdata", "textfmt", "-i="+manager.mergedCoverageDir, "-o="+manager.coverageProfile)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate text report: %w", err)
	}
	return nil
}
