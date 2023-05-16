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
func GetStackReport(
	ctx context.Context,
	manifest *Manifest,
	manifestPath string,
	tokens map[backend.HeaderTokenKeyType]string,
	jsonOut, verboseOut bool,
) error {
	prompts.TelemetryConsentSelect() // if telemetry consent is not set, ask for it
	// prepare telemetry track event properties
	telemetry.SetProperty(ctx, telemetry.KeyJSonOutput, jsonOut)
	telemetry.SetProperty(ctx, telemetry.KeyVerboseOutput, verboseOut)
	telemetry.SetProperty(ctx, telemetry.KeyManifest, manifest.Filename)
	telemetry.SetProperty(ctx, telemetry.KeyEcosystem, manifest.Ecosystem)
	// handle tokens to be included as request headers to the backend
	if _, ok := tokens[backend.HeaderTokenSnyk]; ok {
		telemetry.SetProperty(ctx, telemetry.KeySnykTokenAssociated, true)
	} else if !jsonOut {
		fmt.Println("consider configuring a snyk token using `crda config` to include private snyk vulnerabilities in your report, https://app.snyk.io/redhat/snyk-token")
	}
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
	backendHost := viper.GetString(config.KeyBackendHost.ToString())
	// get stack report response from backend
	response, err := backend.AnalyzeDependencyTree(
		backendHost,
		manifest.Ecosystem,
		cliClient,
		contentType,
		content,
		tokens,
		jsonOut,
	)
	if err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}
	// analyse the response and print the summary
	if err := parseResponse(response, manifest, verboseOut); err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}

	telemetry.SetProperty(ctx, telemetry.KeySuccess, true)
	return nil
}

// parseResponse is used to deserialize a dependency analytics http response
// handles application/json and multipart/mixed(application/json + text/html) response types
// will return error if parsing failed or unknown response body type was used
func parseResponse(response *http.Response, manifest *Manifest, verboseOut bool) error {
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("analyze dependencies request failed, %s", response.Status)
	}

	bodyType, params, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	switch bodyType {
	case "application/json":
		return handleJsonResponse(response.Body, verboseOut)
	case "multipart/mixed":
		return handleMixedResponse(response.Body, params, manifest.Ecosystem, verboseOut)
	default:
		return fmt.Errorf("content type %s is not supported", bodyType)
	}
}
