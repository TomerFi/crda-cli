package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var validShells = []string{"bash", "zsh", "fish", "powershell"}

var completionCmd = &cobra.Command{
	Use:   fmt.Sprintf("completion %s", strings.Join(validShells, "|")),
	Short: "Generate a completions script",
	Long: `Generate a completion script for your shell:

	crda completion bash
	crda completion zsh
	crda completion fish
	crda completion powershell
	`,
	DisableFlagsInUseLine: true, // flags such as debug are useless here, will mess up the script
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	ValidArgs: validShells,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:      printCompletionScript,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

// printCompletionScript is used as the main function for the completion command
// prints a completion script for the selected shell in arg0
// returns error when failed to generate the script or shell unknown
func printCompletionScript(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletion(os.Stdout)
	default:
		return fmt.Errorf("unknown shell %s", args[0])
	}
}
