package coverage

import (
	"fmt"
	"os"
	"path/filepath"
)

type Runner struct {
	coverageDir      string
	coverageFileName string
}

func NewRunner(coverageDir string, coverageFileName string) *Runner {
	return &Runner{
		coverageDir:      coverageDir,
		coverageFileName: coverageFileName,
	}
}

func (runner *Runner) RunCoverage() error {
	if err := os.MkdirAll(runner.coverageDir, 0755); err != nil {
		return fmt.Errorf("failed to create coverage directory: %w", err)
	}

	modules, err := findGoModules(".")
	if err != nil {
		return fmt.Errorf("failed to find go modules: %w", err)
	}
	fmt.Println(len(modules))

	var coverages [][]byte
	for _, module := range modules {
		if err := module.generateCoverage(); err != nil {
			return fmt.Errorf("failed to generate coverage in %s: %w", module.name, err)
		}
		if err := module.saveCoverageToFile(runner.coverageDir); err != nil {
			return fmt.Errorf("failed to write coverage to file: %w", err)
		}
		coverages = append(coverages, module.coverageData)
	}

	mergedCoverage, err := mergeCoverages(coverages)
	if err != nil {
		return fmt.Errorf("failed to merge coverage contents: %w", err)
	}

	mergedCoveragePath := filepath.Join(runner.coverageDir, runner.coverageFileName)
	if err := os.WriteFile(mergedCoveragePath, mergedCoverage, 0644); err != nil {
		return fmt.Errorf("failed to write merged coverage file: %w", err)
	}
	return nil
}
