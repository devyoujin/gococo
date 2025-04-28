package report

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

const (
	htmlReportFileName = "coverage.html"
	textReportFileName = "coverage.txt"
)

type reporter struct {
	coverageDir string
	coveragePath string
}

func NewReporter (coverageDir string, coveragePath string) *reporter {
	return &reporter{
		coverageDir: coverageDir,
		coveragePath: coveragePath,
	}
}

func (reporter *reporter) GenerateHtmlReport() error {
	outputPath := filepath.Join(reporter.coverageDir, htmlReportFileName)
	cmd := exec.Command("go", "tool", "cover", "-html="+reporter.coveragePath, "-o", outputPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate html report: %w", err)
	}
	return nil
}

func (reporter *reporter) GenerateTextReport() error {
	outputPath := filepath.Join(reporter.coverageDir, textReportFileName)
	cmd := exec.Command("go", "tool", "cover", "-func="+reporter.coveragePath, "-o", outputPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate text report: %w", err)
	}
	return nil
}