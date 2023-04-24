package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/segmentio/analytics-go"
	"golang.org/x/exp/slices"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/spf13/cobra"
)

var (
	debug   bool
	client  string
	noColor bool
)

var rootCmd = cobra.Command{
	Use:   "crda",
	Short: "CLI for interacting with the Crda platform",
	Long:  "Use this tool for CodeReady Dependency Analytics reports",
	// error handling is done by the handleErrors function
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// init is used for initializing, parsing, and verifying the global flags
func init() {
	// global flags
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Set DEBUG log level")
	rootCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "c", false, "Toggle colors in output.")
	rootCmd.PersistentFlags().StringVarP(&client, "client", "m", "terminal", "The invoking client for telemetry")
	// parse the flags manually before executing the root command
	if err := rootCmd.ParseFlags(os.Args); err != nil {
		if !(strings.HasPrefix(err.Error(), "unknown flag")) {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	// verify the client arg
	supportedClients := []string{"jenkins", "terminal", "tekton", "gh-actions", "intellij", "vscode", "image"}
	if !slices.Contains(supportedClients, client) {
		fmt.Printf("supported clients are %s", strings.Join(supportedClients, ", "))
		fmt.Println()
		os.Exit(1)
	}

	color.NoColor = noColor // no colors in output, useful for ci
}

// Run is used to initialize telemetry, config, logger, and execute the command
//
// exit code 0 (no error, no vulnerabilities)
// exit code 1 (error)
// exit code 2 (no error, found vulnerabilities)
func Run(segmentClient analytics.Client, userIdFile, configDirectory string) int {
	startTime := time.Now() // set start time for duration calculation

	utils.ConfigureLogging(debug) // configure logging with the debug flag

	// load crda config from env vars or $HOME/.crda/config.yaml
	if err := config.Load(configDirectory); err != nil {
		fmt.Println(err.Error())
		return 1
	}

	// the initiated context includes a value for collecting telemetry properties
	ctx := telemetry.GetContext(context.Background())
	telemetry.SetProperty(ctx, telemetry.KeyClient, client)

	utils.Logger.Debug("executing root command")
	cmd, err := rootCmd.ExecuteContextC(ctx)

	// handle errors and save the exit code (will be returned eventually)
	exitCode := handleErrors(cmd.Context(), err, cmd.Usage)
	telemetry.SetProperty(ctx, telemetry.KeyExitCode, exitCode)

	// get or create telemetry user id and push track event
	if userId, err := telemetry.GetCreateUserIdentity(userIdFile); err != nil {
		utils.Logger.Debugf("no user id to push telemetry to segment, %e", err)
	} else {
		if err := telemetry.PushEvent(ctx, segmentClient, userId, cmd.CommandPath(), startTime); err != nil {
			utils.Logger.Debug("failed to push telemetry event, %e", err)
		}
	}

	utils.Logger.Debugf("exiting with code %d", exitCode)
	return exitCode
}

// handleErrors will analyze the error and will return
// 1 and invoke the usage function if error found
// 2 if no error found but found vulnerabilities
// 0 if no error and no vulnerabilities found
func handleErrors(ctx context.Context, err error, usage func() error) int {
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println()
		_ = usage()
		return 1
	} else if val, ok := telemetry.GetProperty(ctx, telemetry.KeyTotalVulnerabilities); ok {
		if vulnerabilities, _ := strconv.Atoi(val); vulnerabilities > 0 {
			return 2
		}
	}
	return 0
}
