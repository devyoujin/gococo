package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devyoujin/gococo/internal/coverage"
	"github.com/devyoujin/gococo/internal/report"
	"github.com/spf13/cobra"
)

const (
	defaultCoverageDirectory = ".gococo"
	defaultCoverageFileName = "coverage.out"
)

func main() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute command: %+v", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gococo",
	Short: "A CLI tool to generate consolidated test coverage for Go multi-module projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		coverageRunner := coverage.NewRunner(defaultCoverageDirectory, defaultCoverageFileName)
		if err := coverageRunner.RunCoverage(); err != nil {
			return fmt.Errorf("failed to run coverage: %w", err)
		}
		reporter := report.NewReporter(defaultCoverageDirectory, filepath.Join(defaultCoverageDirectory, defaultCoverageFileName))
		if err := reporter.GenerateTextReport(); err != nil {
			return fmt.Errorf("failed to generate text report: %w", err)
		}
		if err := reporter.GenerateHtmlReport(); err != nil {
			return fmt.Errorf("failed to generate html report: %w", err)
		}
		return nil
	},
}

func NewRootCommand() *cobra.Command {
	return rootCmd
}

