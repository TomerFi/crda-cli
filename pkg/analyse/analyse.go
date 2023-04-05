package analyse

import (
	"context"
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/prompts"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
)

func StackReport(ctx context.Context, manifest *Manifest, manifestPath string, json, verbose bool) error {
	prompts.TelemetryConsentSelect()

	telemetry.SetProperty(ctx, telemetry.KeyJSonOutput, json)
	telemetry.SetProperty(ctx, telemetry.KeyVerboseOutput, verbose)
	telemetry.SetProperty(ctx, telemetry.KeyManifest, manifest.Filename)
	telemetry.SetProperty(ctx, telemetry.KeyEcosystem, manifest.Ecosystem)

	// TODO remove this once support for go, node_js, and python_pip is done
	if manifest.Analyzer == nil {
		return fmt.Errorf("sorry, this is a wip, support for %s is not yet active", manifest.Filename)
	}

	if err := manifest.Analyze(ctx, manifest.Ecosystem, manifestPath, json, verbose); err != nil {
		telemetry.SetProperty(ctx, telemetry.KeySuccess, false)
		telemetry.SetProperty(ctx, telemetry.KeyError, telemetry.MaskErrorContent(err))
		return err
	}

	telemetry.SetProperty(ctx, telemetry.KeySuccess, true)
	return nil
}
