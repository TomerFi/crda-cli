package auth

import (
	"context"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/prompts"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
)

func AuthenticateUser(ctx context.Context, snykToken string) error {
	prompts.TelemetryConsentSelect() // if telemetry consent is not set, ask for it
	// collect data required for sending requests to the backend
	cliClient, _ := telemetry.GetProperty(ctx, telemetry.KeyClient)
	host := viper.GetString(config.KeyOldHost.ToString())                   // TODO remove this once done with old backend
	threeScaleToken := viper.GetString(config.KeyOld3ScaleToken.ToString()) // TODO remove this once done with old backend

	// if we don't already have a crda user key, ask the backend for an ew one
	var crdaKey string
	if !viper.IsSet(config.KeyCrdaKey.ToString()) {
		if newUserKey, err := backend.RequestNewUserKey(host, threeScaleToken, cliClient); err == nil {
			crdaKey = newUserKey
			viper.Set(config.KeyCrdaKey.ToString(), newUserKey)
		}
	} else {
		crdaKey = viper.GetString(config.KeyCrdaKey.ToString())
	}

	// if snyk token doesn't match regex, prompt the user to input token
	token := snykToken
	if !utils.MatchSnykRegex(token) { // NOTE this will return false for empty strings
		if tok, err := prompts.SnykTokenPrompt(); err != nil {
			return err
		} else {
			token = tok // NOTE this might be an empty string if not token was input by the user
		}
	}

	// if we have a user id and a token, associated the token with the crda key
	tokenAssociated := false
	if crdaKey != "" && token != "" {
		if err := backend.AssociateSnykToken(host, threeScaleToken, cliClient, crdaKey, token); err != nil {
			return err
		} else {
			tokenAssociated = true
		}
	} else {
		utils.Logger.Debug("no crda or token, association skipped")
	}
	telemetry.SetProperty(ctx, telemetry.KeySnykTokenAssociated, tokenAssociated)

	return nil
}
