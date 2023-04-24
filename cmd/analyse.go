package cmd

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"os"
	"strings"

	"github.com/rhecosystemappeng/crda-cli/pkg/analyse"
	"github.com/spf13/cobra"
)

var (
	jsonOutput    bool
	verboseOutput bool
)

var analyseCmd = &cobra.Command{
	Use:   fmt.Sprintf("analyse %s", strings.Join(analyse.SupportedManifestsFilenames, "|")),
	Short: "Preform dependency analysis report",
	Long:  "Preform dependency analysis report. Will exit with status code 2 if vulnerabilities found",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), isSupportedPath),
	RunE: printReport,
}

// init is used for setting the flags and binding the command
func init() {
	analyseCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Get report in a JSON format")
	analyseCmd.Flags().BoolVarP(&verboseOutput, "verbose", "v", false, "Get detailed report")
	rootCmd.AddCommand(analyseCmd)
}

// printReport will print the requested report (stack analysis)
// returns error if failed generating invoking backend analysis
func printReport(cmd *cobra.Command, args []string) error {
	utils.Logger.Debug("executing analyse command")
	file, err := os.Stat(args[0])
	if err != nil {
		return err
	}
	manifest, err := analyse.GetManifest(file.Name())
	if err != nil {
		return err
	}
	return analyse.StackReport(cmd.Context(), manifest, args[0], jsonOutput, verboseOutput)
}

// isSupportedPath will return an error if the manifest file is unsupported/unknown
func isSupportedPath(cmd *cobra.Command, args []string) error {
	return analyse.IsSupportedManifestPath(args[0])
}
