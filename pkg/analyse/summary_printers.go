package analyse

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
)

func printSummary(vulSummary VulnerabilitiesSummary, reportUri string) error {
	// TODO waiting for this https://github.com/RHEcosystemAppEng/crda-backend/issues/28
	// include volSummary in the summary print
	white := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	fmt.Println(white("Full Report: "), reportUri)

	return nil
}

func printJson(vulSummary VulnerabilitiesSummary) error {

	output, err := json.MarshalIndent(vulSummary, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
