package cmd

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get binary version",
	Long:  "Command to output version of the binary",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.NoArgs,
	Run:  printVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// printVersion is used as the main function for the version command
// prints the info for the current version
func printVersion(cmd *cobra.Command, args []string) {
	utils.Logger.Debug("executing version command")
	fmt.Println(utils.BuildVersion())
}
