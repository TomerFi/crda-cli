package analyse

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend/api"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	colorBlueHi    = color.New(color.FgHiBlue, color.Bold)
	colorGreenHi   = color.New(color.FgGreen, color.Bold)
	colorMagentaHi = color.New(color.FgHiMagenta, color.Bold)
	colorRedHi     = color.New(color.FgHiRed, color.Bold)
	colorRed       = color.New(color.FgRed, color.Bold)
	colorWhiteHi   = color.New(color.FgHiWhite, color.Bold)
	colorWhite     = color.New(color.FgWhite, color.Bold)
	colorYellow    = color.New(color.FgHiYellow, color.Bold)
)

func printSummary(report *api.AnalysisReport, reportUri string) {
	fmt.Println("Summary Report for Dependency Analysis")
	fmt.Println()
	colorWhiteHi.PrintlnFunc()("Total Scanned Dependencies:", *report.Summary.Dependencies.Scanned)
	colorWhiteHi.PrintlnFunc()("Total Scanned Transitive Dependencies:", *report.Summary.Dependencies.Transitive)
	colorWhiteHi.PrintlnFunc()("Direct Vulnerable Dependencies:", *report.Summary.Vulnerabilities.Direct)
	colorWhiteHi.PrintlnFunc()("Total Vulnerabilities:", *report.Summary.Vulnerabilities.Total)
	colorRedHi.PrintlnFunc()("Critical Vulnerabilities:", *report.Summary.Vulnerabilities.Critical)
	colorMagentaHi.PrintlnFunc()("High Vulnerabilities:", *report.Summary.Vulnerabilities.High)
	colorYellow.PrintlnFunc()("Medium Vulnerabilities:", *report.Summary.Vulnerabilities.Medium)
	colorBlueHi.PrintlnFunc()("Low Vulnerabilities:", *report.Summary.Vulnerabilities.Low)
	fmt.Println()
	colorWhiteHi.PrintlnFunc()("Full Report:", reportUri)
}

func printVerboseSummary(report *api.AnalysisReport, reportUri string) {
	fmt.Println("Verbose Report for Dependency Analysis:")
	fmt.Println()

	fmt.Println(
		colorWhite.Sprintf(
			"Scanned %g Dependencies and %g Transitives,",
			*report.Summary.Dependencies.Scanned,
			*report.Summary.Dependencies.Transitive,
		),
		colorRed.Sprintf("Found %g Issues", *report.Summary.Vulnerabilities.Total),
	)
	fmt.Println()
	for _, dep := range *report.Dependencies {
		fmt.Println(colorWhiteHi.Sprintf("Direct dependency: %s@%s", *dep.Ref.Name, *dep.Ref.Version))
		if dep.HighestVulnerability != nil {
			fmt.Println(
				colorWhite.Sprint("\t\u2212 ", "Highest vulnerability:"),
				describeVulnerability(dep.HighestVulnerability),
			)
		}
		if len(*dep.Issues) > 0 {
			fmt.Println(colorWhiteHi.Sprintf("\t\u2212 %d Direct Vulnerabilities:", len(*dep.Issues)))
			for _, issue := range *dep.Issues {
				fmt.Println("\t\t\u2718 ", describeVulnerability(&issue))
			}
		}
		if len(*dep.Transitive) > 0 {
			fmt.Println(colorWhiteHi.Sprintf("\t\u2212 %d Transitive Vulnerabilities:", len(*dep.Transitive)))
			for _, transitive := range *dep.Transitive {
				if len(*transitive.Issues) > 0 {
					fmt.Println(colorWhite.Sprint("\t\t\u2212 Transitive dependency:"), colorWhiteHi.Sprintf("%s@%s", *transitive.Ref.Name, *transitive.Ref.Version))
					for _, issue := range *transitive.Issues {
						fmt.Println("\t\t\t\t\u2718 ", describeVulnerability(&issue))
					}
				}
				if len(*transitive.Remediations) > 0 {
					fmt.Println(colorWhiteHi.Sprintf("\t\u2212 %d Remedying actions for transitive found:", len(*transitive.Remediations)))
					for remediedCve, remediation := range *transitive.Remediations {
						colorGreenHi.PrintlnFunc()(
							"\t\t\u2713 ",
							remediedCve,
							"can be remedied with",
							fmt.Sprintf("%s@%s", *remediation.MavenPackage.Name, *remediation.MavenPackage.Version))
					}
				}
			}
		}

		if len(*dep.Remediations) > 0 {
			fmt.Println(colorWhiteHi.Sprintf("\t\u2212 %d Remedying actions found:", len(*dep.Remediations)))
			for remediedCve, remediation := range *dep.Remediations {
				colorGreenHi.PrintlnFunc()(
					"\t\t\u2713 ",
					remediedCve,
					"can be remedied with",
					fmt.Sprintf("%s@%s", *remediation.MavenPackage.Name, *remediation.MavenPackage.Version))
			}
		}
	}
	fmt.Println()
	colorWhiteHi.PrintlnFunc()("Full Report:", reportUri)
}

func printJson(report *api.AnalysisReport, verboseOut bool) error {
	var output []byte
	var err error
	if verboseOut {
		// if verbose output requested print the entire report
		output, err = json.MarshalIndent(report, "", "\t")
	} else {
		// if verbose output NOT requested, print only the summary object
		output, err = json.MarshalIndent(report.Summary, "", "\t")
	}
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func describeVulnerability(issue *api.Issue) string {
	colorFunc := colorWhiteHi
	switch *issue.Severity {
	case api.CRITICAL:
		colorFunc = colorRed
	case api.HIGH:
		colorFunc = colorMagentaHi
	case api.MEDIUM:
		colorFunc = colorYellow
	case api.LOW:
		colorFunc = colorBlueHi
	}

	return colorFunc.Sprintf(
		"%s [%s] %s (%s)",
		*issue.Title,
		cases.Title(language.Und).String(string(*issue.Severity)),
		*issue.Cves,
		utils.GetProviderUrl(*issue.Source, *issue.Id),
	)
}
