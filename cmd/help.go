package cmd

import "github.com/spf13/cobra"

var helpCmd = &cobra.Command{
	Use:    "help",
	Short:  "Help about any command",
	Long:   "Get help for any command in the available commands list",
	Hidden: true,
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.Help()
	},
}

func init() {
	rootCmd.SetHelpCommand(helpCmd)
}
