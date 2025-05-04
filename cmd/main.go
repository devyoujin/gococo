package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	coverageDir = ".gococo"
	mergedCoverageDir = "data"
	coverageProfile  = "coverage.out"
	coverageReportHtml = "coverage.html"
	coverageReportText = "coverage.txt"
)

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gococo",
		Short: "A CLI tool to generate consolidated test coverage for Go multi-module projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			coverageRunner := NewRunner(coverageDir, mergedCoverageDir, coverageProfile, coverageReportHtml, coverageReportText)
			if err := coverageRunner.Run(); err != nil {
				return fmt.Errorf("failed to run coverage: %w", err)
			}
			return nil
		},
	}
}
