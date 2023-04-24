package cmd

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"strings"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage crda config",
	Long:  "Command used for managing crda config",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
}

var configGetCmd = &cobra.Command{
	Use:   fmt.Sprintf("get [%s]", strings.Join(config.KnownConfigKeyStrings, "|")),
	Short: "Get crda config",
	Long:  "Display crda config",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.MatchAll(cobra.MaximumNArgs(1), verifyKeyExist),
	Run:  getFromConfig,
}

var configSetCmd = &cobra.Command{
	Use:   fmt.Sprintf("set {%s your-value}", strings.Join(config.KnownConfigKeyStrings, "|")),
	Short: "Set crda config key.s",
	Long:  "Set a crda config key",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.MatchAll(cobra.ExactArgs(2), verifyKeyValueArgs),
	RunE: setConfigKeyValue,
}

// init is used to bind the get/set commands to the config command
// and bind the config command to the root command
func init() {
	configCmd.AddCommand(configGetCmd, configSetCmd)
	rootCmd.AddCommand(configCmd)
}

// verifyKeyExist is used as a validArgs function for the config get command
// returns error if key in arg0 doesn't exist in global config
// if no arg0, returns nil to allow "get all"
func verifyKeyExist(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && !viper.IsSet(args[0]) {
		return fmt.Errorf("config key %s is not set", args[0])
	}
	return nil
}

// verifyKeyValueArgs is used as a validArgs function for the config set command
// returns error if key in arg0 is not a known config key
func verifyKeyValueArgs(cmd *cobra.Command, args []string) error {
	if !slices.Contains(config.KnownConfigKeyStrings, args[0]) {
		return fmt.Errorf("supported config keys are %s", strings.Join(config.KnownConfigKeyStrings, ", "))
	}
	return nil
}

// getFromConfig is used as the main function for the config get command
// prints a specific 'key: value' or all configuration
// key is expected in arg0, if not found, will print all config pairs
func getFromConfig(cmd *cobra.Command, args []string) {
	utils.Logger.Debug("executing config get command")
	if len(args) > 0 {
		fmt.Printf("%s: %s", args[0], viper.GetString(args[0]))
		fmt.Println()
	} else {
		for k, v := range viper.AllSettings() {
			fmt.Printf("%s: %s", k, v)
			fmt.Println()
		}
	}
}

// setConfigKeyValue is used as the main function for the config set command
// key is expected in arg0 and value is in arg1
// saved to the global config file and printed to the console
// will return error when failed writing the config file
func setConfigKeyValue(cmd *cobra.Command, args []string) error {
	viper.Set(args[0], args[1])
	fmt.Printf("%s: %s", args[0], args[1])
	fmt.Println()
	return viper.WriteConfig()
}
