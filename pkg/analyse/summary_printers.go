package analyse

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend/api"
)

var (
	yellow  = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	white   = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	red     = color.New(color.FgHiRed, color.Bold).SprintFunc()
	blue    = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	magenta = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
)

func printSummary(report *api.AnalysisReport, reportUri string) {
	fmt.Println("Summary Report for Dependency Analysis:")
	fmt.Println()
	fmt.Println(white("Total Scanned Dependencies: "), *report.Summary.Dependencies.Scanned)
	fmt.Println(white("Total Scanned Transitive Dependencies: "), white(*report.Summary.Dependencies.Transitive))
	fmt.Println(white("Direct Vulnerable Dependencies: "), white(*report.Summary.Vulnerabilities.Direct))
	fmt.Println(white("Total Vulnerabilities: "), white(*report.Summary.Vulnerabilities.Total))
	// TODO missing from the new backend design: "Total Direct Vulnerabilities Ignored"
	// TODO missing from the new backend design: "Total Transitive Vulnerabilities Ignored"
	// TODO missing from the new backend implementation: "Publicly Available Vulnerabilities"
	// TODO missing from the new backend implementation: "Snyk Token Registered"
	// TODO missing from the new backend implementation: "Vulnerabilities Unique to Snyk"
	fmt.Println(red("Critical Vulnerabilities: "), red(*report.Summary.Vulnerabilities.Critical))
	fmt.Println(magenta("High Vulnerabilities: "), magenta(*report.Summary.Vulnerabilities.High))
	fmt.Println(yellow("Medium Vulnerabilities: "), yellow(*report.Summary.Vulnerabilities.Medium))
	fmt.Println(blue("Low Vulnerabilities: "), blue(*report.Summary.Vulnerabilities.Low))
	fmt.Println()
	fmt.Println(white("Full Report: "), reportUri)
}

func printJson(report *api.AnalysisReport) error {
	output, err := json.MarshalIndent(report.Summary, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
