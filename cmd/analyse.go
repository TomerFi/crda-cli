package cmd

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/rhecosystemappeng/crda-cli/pkg/analyse"
	"github.com/spf13/cobra"
)

var (
	jsonOutput    bool
	verboseOutput bool
	snykToken     string
)

var analyseCmd = &cobra.Command{
	Use:   fmt.Sprintf("analyse %s", strings.Join(analyse.SupportedManifestsFilenames, "|")),
	Short: "Preform dependency analysis report",
	Long:  "Preform dependency analysis report. Will exit with status code 2 if vulnerabilities found",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), isSupportedPath),
	RunE: analyseManifest,
}

// init is used for setting the flags and binding the command
func init() {
	analyseCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Get report in a JSON format")
	analyseCmd.Flags().BoolVarP(&verboseOutput, "verbose", "v", false, "Get detailed report")
	analyseCmd.Flags().StringVar(&snykToken, "snyk-token", "", "Include snyk token for private snyk vulnerabilities")
	rootCmd.AddCommand(analyseCmd)
}

// analyseManifest will print the requested report (stack analysis)
// returns error if failed generating invoking backend analysis
func analyseManifest(cmd *cobra.Command, args []string) error {
	utils.Logger.Debug("executing analyse command")
	file, err := os.Stat(args[0])
	if err != nil {
		return err
	}
	manifest, err := analyse.GetManifest(file.Name())
	if err != nil {
		return err
	}
	// TODO remove this once implement support for go, node_js, and python_pip is done
	if manifest.TreeProvider == nil {
		return fmt.Errorf("sorry, this is a wip, support for %s is not yet active", manifest.Filename)
	}
	// collect tokens to be included as request headers to the backend
	var tokens = make(map[backend.HeaderTokenKeyType]string)
	if utils.MatchSnykRegex(snykToken) {
		utils.Logger.Debug("using snyk token from flag")
		tokens[backend.HeaderTokenSnyk] = snykToken
	} else if viper.IsSet(config.KeySnykToken.ToString()) {
		utils.Logger.Debug("using snyk token from env/config")
		tokens[backend.HeaderTokenSnyk] = viper.GetString(config.KeySnykToken.ToString())
	}

	return analyse.GetStackReport(cmd.Context(), manifest, args[0], tokens, jsonOutput, verboseOutput)
}

// isSupportedPath will return an error if the manifest file is unsupported/unknown
func isSupportedPath(cmd *cobra.Command, args []string) error {
	return analyse.IsSupportedManifestPath(args[0])
}
