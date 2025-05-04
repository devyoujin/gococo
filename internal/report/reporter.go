package report

import (
	"fmt"
	"os/exec"

	"github.com/devyoujin/gococo/internal/utils"
)

type Reporter struct {
	executor           utils.CommandExecutor
	coverageProfile    string
	htmlReportFileName string
	textReportFileName string
}

func NewReporter(
	executor utils.CommandExecutor,
	coverageProfile string,
	htmlReportFileName string,
	textReportFileName string) *Reporter {
	return &Reporter{
		executor:           executor,
		coverageProfile:    coverageProfile,
		htmlReportFileName: htmlReportFileName,
		textReportFileName: textReportFileName,
	}
}

func (reporter *Reporter) GenerateHtmlReport() error {
	cmd := exec.Command("go", "tool", "cover", "-html="+reporter.coverageProfile, "-o", reporter.htmlReportFileName)
	err := reporter.executor.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to generate html report: %w", err)
	}
	return nil
}

func (reporter *Reporter) GenerateTextReport() error {
	cmd := exec.Command("go", "tool", "cover", "-func="+reporter.coverageProfile, "-o="+reporter.textReportFileName)
	err := reporter.executor.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to generate text report: %w", err)
	}
	return nil
}
