package cmd

import (
	"github.com/rhecosystemappeng/crda-cli/pkg/auth"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var snykToken string

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Link crda user with snyk",
	Long:  "Link crda user id with provider token, i.e. Snyk to unlock Verbose stack analyses",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.NoArgs,
	RunE: authenticateWithToken,
}

func init() {
	authCmd.Flags().StringVarP(&snykToken, "snyk-token", "t", "", "Token for Snyk Authentication")
	rootCmd.AddCommand(authCmd)
}

func authenticateWithToken(cmd *cobra.Command, _ []string) error {
	utils.Logger.Debug("executing auth command")
	return auth.AuthenticateUser(cmd.Context(), snykToken)
}
