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

// printSummary is used for printing a simple summary of the report
// takes a reference to the deserialized JSON Analysis report and an uri for the generated html one
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

// printVerboseSummary is used for printing a detailed summary of the report
// takes a reference to the deserialized JSON Analysis report and an uri for the generated html one
func printVerboseSummary(report *api.AnalysisReport, reportUri string) {
	fmt.Println("Verbose Report for Dependency Analysis:")
	fmt.Println()

	// print summary, i.e. "Scanned 10 Dependencies and 192 Transitives, Found 14 Issues"
	fmt.Println(
		colorWhite.Sprintf(
			"Scanned %g Dependencies and %g Transitives,",
			*report.Summary.Dependencies.Scanned,
			*report.Summary.Dependencies.Transitive,
		),
		colorRed.Sprintf("Found %g Issues", *report.Summary.Vulnerabilities.Total),
	)
	fmt.Println()
	// iterate over all dependencies in the response
	for _, dep := range *report.Dependencies {
		// print direct dependency title, i.e. "Direct dependency: io.quarkus:quarkus-hibernate-orm@2.13.5.Final"
		fmt.Println(colorWhiteHi.Sprintf("Direct dependency: %s@%s", *dep.Ref.Name, *dep.Ref.Version))
		// if highest vulnerability provided print (note the tab prefix), i.e.
		// "        − Highest vulnerability: Information Exposure [High] [CVE-2023-21930] (https://security.snyk.io/vuln/SNYK-JAVA-ORGGRAALVMSDK-5457933)"
		if dep.HighestVulnerability != nil {
			fmt.Println(
				colorWhite.Sprint("\t\u2212 ", "Highest vulnerability:"),
				describeVulnerability(dep.HighestVulnerability),
			)
		}
		if len(*dep.Issues) > 0 {
			// if found direct vulnerabilities print title (note the tab prefix), i.e.
			// "       − 2 Direct Vulnerabilities:"
			fmt.Println(colorWhiteHi.Sprintf("\t\u2212 %d Direct Vulnerabilities:", len(*dep.Issues)))
			// print summary of every direct issue found (note the double tab prefix), i.e.
			// "                ✘  Access Restriction Bypass [High] [CVE-2022-4147] (https://security.snyk.io/vuln/SNYK-JAVA-IOQUARKUS-3149918)"
			for _, issue := range *dep.Issues {
				fmt.Println("\t\t\u2718 ", describeVulnerability(&issue))
			}
		}
		if len(*dep.Transitive) > 0 {
			// if found transitive vulnerabilities print title (note the tab prefix), i.e.
			// "        − 3 Transitive Vulnerabilities:"
			fmt.Println(colorWhiteHi.Sprintf("\t\u2212 %d Transitive Vulnerabilities:", len(*dep.Transitive)))
			for _, transitive := range *dep.Transitive {
				if len(*transitive.Issues) > 0 {
					// print the transitive dependency name (note the double tab prefix), i.e.
					// "                − Transitive dependency: io.netty:netty-handler@4.1.82.Final"
					fmt.Println(colorWhite.Sprint("\t\t\u2212 Transitive dependency:"), colorWhiteHi.Sprintf("%s@%s", *transitive.Ref.Name, *transitive.Ref.Version))
					// print summary of every transitive issue found (note the triple tab prefix), i.e.
					// "                        ✘  Improper Certificate Validation [Medium] [] (https://security.snyk.io/vuln/SNYK-JAVA-IONETTY-1042268)"
					for _, issue := range *transitive.Issues {
						fmt.Println("\t\t\t\u2718 ", describeVulnerability(&issue))
					}
				}
				if len(*transitive.Remediations) > 0 {
					// if found remediating actions for transitive vulnerabilities print title (note the double tab prefix), i.e.
					// "                − 1 Remedying actions found:"
					fmt.Println(colorWhiteHi.Sprintf("\t\t\u2212 %d Remedying actions for transitive found:", len(*transitive.Remediations)))
					// print summary of every remediating action for transitive issue found (note the triple tab prefix), i.e.
					// "                       ✓  CVE-2023-0044 can be remedied with io.quarkus:quarkus-vertx-http@2.13.7.Final-redhat-00003"
					for remediedCve, remediation := range *transitive.Remediations {
						colorGreenHi.PrintlnFunc()(
							"\t\t\t\u2713 ",
							remediedCve,
							"can be remedied with",
							fmt.Sprintf("%s@%s", *remediation.MavenPackage.Name, *remediation.MavenPackage.Version))
					}
				}
			}
		}

		if len(*dep.Remediations) > 0 {
			// if found remediating actions for direct vulnerabilities print title (note the tab prefix), i.e.
			// "        − 1 Remedying actions found:"
			fmt.Println(colorWhiteHi.Sprintf("\t\u2212 %d Remedying actions found:", len(*dep.Remediations)))
			// print summary of every remediating action for transitive issue found (note the double tab prefix), i.e.
			// "               ✓  CVE-2023-0044 can be remedied with io.quarkus:quarkus-vertx-http@2.13.7.Final-redhat-00003"
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
	// print generated html report uri
	colorWhiteHi.PrintlnFunc()("Full Report:", reportUri)
}

// printJson is used for printing the entire JSON response to the standard output
// takes a reference to the deserialized JSON Analysis report and a verbose boolean
// will print the entire response if verbose is true or just the summary object if not
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

// describeVulnerability is used to turn a vulnerability issue to a one-line string
// takes an issue reference and return a pattern of "TITLE [SEVERITY] [CVE LIST] (PROVIDER URL)"
// i.e. "Directory Traversal [Medium] [CVE-2023-24815] (https://security.snyk.io/vuln/SNYK-JAVA-IOVERTX-3318108)"
func describeVulnerability(issue *api.Issue) string {
	// choose color function based on severity
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
