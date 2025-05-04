package coverage

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"

	"github.com/devyoujin/gococo/internal/utils"
)

type ManagerInterface interface {
	FindGoModules() ([]module, error)
	GenerateCoverages(modules []module) error
	GenerateCoverProfile() error
}

type Manager struct {
	executor          utils.CommandExecutor
	pwd               string
	mergedCoverageDir string
	coverageProfile   string
}

func NewManager(
	executor utils.CommandExecutor,
	pwd string,
	mergedCoverageDir string,
	coverageProfile string) ManagerInterface {
	return &Manager{
		pwd:               pwd,
		mergedCoverageDir: mergedCoverageDir,
		coverageProfile:   coverageProfile,
		executor:          executor,
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
		cmd := exec.Command("go", "test", "-cover", "./...", "-test.gocoverdir="+manager.mergedCoverageDir)
		cmd.Dir = module.path
		if err := manager.executor.Run(cmd); err != nil {
			return fmt.Errorf("failed to run tests in %s: %w", module.path, err)
		}
	}
	return nil
}

func (manager *Manager) GenerateCoverProfile() error {
	cmd := exec.Command("go", "tool", "covdata", "textfmt", "-i="+manager.mergedCoverageDir, "-o="+manager.coverageProfile)
	if err := manager.executor.Run(cmd); err != nil {
		return fmt.Errorf("failed to generate coverage profile: %w", err)
	}
	return nil
}
