package analyse

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/prompts"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/spf13/viper"
	"mime"
	"mime/multipart"
	"net/http"
)

type summaryAndReport struct {
	VulnerabilitiesSummary
	ReportFileUri string
}

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
	)
	if err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}
	// analyse the response and get the vulnerabilities summary and report uri
	sumAndReport, err := parseResponse(response, manifest)
	if err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}
	// print the report uri and summary to the standard output
	if err := printSummary(sumAndReport, jsonOut); err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}

	telemetry.SetProperty(ctx, telemetry.KeySuccess, true)
	return nil
}

// parseResponse is used to verify the backend stack analysis response, parse its body,
// and return a summary including the local uri for the report
// it will return an error if the response is not ok, it used an unknown body type,
// or failed to deserialize and parse the body
func parseResponse(response *http.Response, manifest *Manifest) (summaryAndReport, error) {
	if response.StatusCode != http.StatusOK {
		return summaryAndReport{}, fmt.Errorf("analyze dependencies request failed, %s", response.Status)
	}

	bodyType, params, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return summaryAndReport{}, err
	}
	if "multipart/mixed" != bodyType {
		return summaryAndReport{}, fmt.Errorf("content type %s is not supported", bodyType)
	}

	var vulSummary VulnerabilitiesSummary
	var reportUri string

	multipartReader := multipart.NewReader(response.Body, params["boundary"])
	for part, err := multipartReader.NextPart(); err == nil; part, err = multipartReader.NextPart() {
		partType := part.Header.Get("Content-Type")
		switch partType {
		case "application/json":
			if reports, err := backend.ParseJsonResponse(part); err != nil {
				return summaryAndReport{}, err
			} else {
				vulSummary = processVulnerabilities(reports)
			}
		case "text/html":
			if reportUri, err = backend.ParseHtmlResponse(part, manifest.Ecosystem); err != nil {
				return summaryAndReport{}, err
			}
		default:
			return summaryAndReport{}, fmt.Errorf("unknown response type %s", partType)
		}
	}

	return summaryAndReport{vulSummary, reportUri}, nil
}

// printSummary is used to print the summary to the standard output as a human-readable data
// use jsonOut=true to print as a machine-readable json object
func printSummary(sumAndReport summaryAndReport, jsonOut bool) error {
	if jsonOut {
		output, err := json.MarshalIndent(sumAndReport, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	} else {
		// TODO waiting for this https://github.com/RHEcosystemAppEng/crda-backend/issues/28
		// include volSummary in the summary print
		white := color.New(color.FgHiWhite, color.Bold).SprintFunc()
		fmt.Println(white("Full Report: "), sumAndReport.ReportFileUri)
	}

	return nil
}
