package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devyoujin/gococo/internal/coverage"
	"github.com/devyoujin/gococo/internal/report"
)

type Runner struct {
	coverageManager   coverage.ManagerInterface
	coverageReporter  *report.Reporter
	coverageDir       string
	mergedCoverageDir string
}

func NewRunner(coverageDir string, mergedCoverageDir string, coverageProfile string, coverageReportHtml string, coverageReportText string) *Runner {
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get current working directory: %w", err))
	}
	coverageDir = filepath.Join(pwd, coverageDir)
	mergedCoverageDir = filepath.Join(coverageDir, mergedCoverageDir)
	coverageProfile = filepath.Join(coverageDir, coverageProfile)
	coverageReportHtml = filepath.Join(coverageDir, coverageReportHtml)
	coverageReportText = filepath.Join(coverageDir, coverageReportText)

	coverageManager := coverage.NewManager(pwd, mergedCoverageDir, coverageProfile)
	coverageReporter := report.NewReporter(coverageProfile, coverageReportHtml, coverageReportText)
	return &Runner{
		coverageManager:   coverageManager,
		coverageReporter:  coverageReporter,
		coverageDir:       coverageDir,
		mergedCoverageDir: mergedCoverageDir,
	}
}

func (runner *Runner) Run() error {
	var err error
	modules, err := runner.coverageManager.FindGoModules()
	if err != nil {
		return fmt.Errorf("failed to find go modules: %w", err)
	}
	if len(modules) == 0 {
		fmt.Print("no go modules found in the current directory")
		return nil
	}
	err = os.RemoveAll(runner.coverageDir)
	if err != nil {
		return fmt.Errorf("failed to remove coverage data directory: %w", err)
	}
	err = os.MkdirAll(runner.mergedCoverageDir, 0755);
	if err != nil {
		return fmt.Errorf("failed to create coverage data directory: %w", err)
	}
	err = runner.coverageManager.GenerateCoverages(modules)
	if err != nil {
		return fmt.Errorf("failed to generate coverages: %w", err)
	}
	err = runner.coverageManager.GenerateCoverProfile()
	if err != nil {
		return fmt.Errorf("failed to generate coverage profile: %w", err)
	}
	err = runner.coverageReporter.GenerateHtmlReport()
	if err != nil {
		return fmt.Errorf("failed to generate html report: %w", err)
	}
	err = runner.coverageReporter.GenerateTextReport()
	if err != nil {
		return fmt.Errorf("failed to generate text report: %w", err)
	}
	return nil
}
