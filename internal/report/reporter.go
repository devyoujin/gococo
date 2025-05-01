package report

import (
	"fmt"
	"os/exec"
)

type Reporter struct {
	coverageProfile string
	htmlReportFileName string
	textReportFileName string
}

func NewReporter(
	coverageProfile string,
	htmlReportFileName string,
	textReportFileName string,
	) *Reporter {
	return &Reporter{
		coverageProfile: coverageProfile,
		htmlReportFileName: htmlReportFileName,
		textReportFileName: textReportFileName,
	}
}

func (reporter *Reporter) GenerateHtmlReport() error {
	cmd := exec.Command("go", "tool", "cover", "-html="+reporter.coverageProfile, "-o", reporter.htmlReportFileName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate html report: %w", err)
	}
	return nil
}

func (reporter *Reporter) GenerateTextReport() error {
	cmd := exec.Command("go", "tool", "cover", "-func="+reporter.coverageProfile, "-o="+reporter.textReportFileName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate text report: %w", err)
	}
	return nil
}
