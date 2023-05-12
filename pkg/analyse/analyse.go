package analyse

import (
	"context"
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/prompts"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/spf13/viper"
	"mime"
	"net/http"
)

// GetStackReport is used for requesting a stack analysis from the backend server
// It will print a human-readable report summary to the standard output
// Use jsonOut=true to print the summary as a machine-readable json object
// Use verbose=true to include private vulnerabilities in the report
func GetStackReport(ctx context.Context, manifest *Manifest, manifestPath string, jsonOut, verbose bool) error {
	prompts.TelemetryConsentSelect() // if telemetry consent is not set, ask for it
	// prepare telemetry track event properties
	telemetry.SetProperty(ctx, telemetry.KeyJSonOutput, jsonOut)
	telemetry.SetProperty(ctx, telemetry.KeyVerboseOutput, verbose)
	telemetry.SetProperty(ctx, telemetry.KeyManifest, manifest.Filename)
	telemetry.SetProperty(ctx, telemetry.KeyEcosystem, manifest.Ecosystem)
	// get the content and content type from the concrete tree provider
	// these will get delegated to the backend
	content, contentType, err := manifest.TreeProvider.Provide(ctx, manifestPath)
	if err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}
	// collect data required for sending requests to the backend
	cliClient, _ := telemetry.GetProperty(ctx, telemetry.KeyClient)
	oldHost := viper.GetString(config.KeyOldHost.ToString())                // TODO remove this once done with old backend
	threeScaleToken := viper.GetString(config.KeyOld3ScaleToken.ToString()) // TODO remove this once done with old backend
	backendHost := viper.GetString(config.KeyBackendHost.ToString())
	// if we don't already have a crda user key, ask the backend for a new one
	var crdaKey string
	if !viper.IsSet(config.KeyCrdaKey.ToString()) {
		if newUserKey, err := backend.RequestNewUserKey(oldHost, threeScaleToken, cliClient); err == nil {
			crdaKey = newUserKey
			viper.Set(config.KeyCrdaKey.ToString(), newUserKey)
		}
	} else {
		crdaKey = viper.GetString(config.KeyCrdaKey.ToString())
	}
	// get stack report response from backend
	response, err := backend.AnalyzeDependencyTree(
		backendHost,
		JavaMaven.Ecosystem,
		crdaKey,
		cliClient,
		contentType,
		content,
		jsonOut,
	)
	if err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}
	// analyse the response and print the summary
	if err := parseResponse(response, manifest); err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}

	telemetry.SetProperty(ctx, telemetry.KeySuccess, true)
	return nil
}

func parseResponse(response *http.Response, manifest *Manifest) error {
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("analyze dependencies request failed, %s", response.Status)
	}

	bodyType, params, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	switch bodyType {
	case "application/json":
		return handleJsonResponse(response.Body)
	case "multipart/mixed":
		return handleMixedResponse(response.Body, params, manifest.Ecosystem)
	default:
		return fmt.Errorf("content type %s is not supported", bodyType)
	}
}
