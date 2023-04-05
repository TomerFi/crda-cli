package analyse

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

type JavaMavenAnalyzer struct{}

func (a *JavaMavenAnalyzer) Analyze(ctx context.Context, ecosystem string, manifestPath string, json, verbose bool) error {
	mvn, err := exec.LookPath("mvn")
	if err != nil {
		return err
	}

	tmpDesTree := filepath.Join(os.TempDir(), "tmp-deps-tree.txt")
	defer os.Remove(tmpDesTree)

	if _, err := os.Stat(tmpDesTree); !errors.Is(err, fs.ErrNotExist) {
		os.Remove(tmpDesTree) // if the temp file exists - remove it
	}

	cleanExec := exec.Command(mvn, "-q", "clean", "-f", manifestPath)
	treeExec := exec.Command(mvn, "-q", "dependency:tree", "-DoutputType=dot", fmt.Sprintf("-DoutputFile=%s", tmpDesTree), "-f", manifestPath)

	if err := cleanExec.Run(); err != nil {
		return err
	}

	if err := treeExec.Run(); err != nil {
		return err
	}

	graph, err := os.ReadFile(tmpDesTree)
	if err != nil {
		return err
	}

	cliClient, _ := telemetry.GetProperty(ctx, telemetry.KeyClient)
	oldHost := viper.GetString(config.KeyOldHost.ToString())                // TODO remove this once done with old backend
	threeScaleToken := viper.GetString(config.KeyOld3ScaleToken.ToString()) // TODO remove this once done with old backend
	backendHost := viper.GetString(config.KeyBackendHost.ToString())

	// if we don't already have a crda user key, ask the backend for an ew one
	var crdaKey string
	if !viper.IsSet(config.KeyCrdaKey.ToString()) {
		if newUserKey, err := backend.RequestNewUserKey(oldHost, threeScaleToken, cliClient); err == nil {
			crdaKey = newUserKey
			viper.Set(config.KeyCrdaKey.ToString(), newUserKey)
		}
	} else {
		crdaKey = viper.GetString(config.KeyCrdaKey.ToString())
	}

	body, err := backend.AnalyzeDependencyTree(
		backendHost,
		crdaKey,
		cliClient,
		"text/vnd.graphviz",
		graph,
	)
	if err != nil {
		return err
	}

	htmlFileUri, err := utils.SaveReportToTempHtmlFile(*body, ecosystem)
	if err != nil {
		return err
	}

	// TODO replace this with logic for printing verbose/non-verbose json/non-json summary
	white := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	fmt.Println(white("Full Report: "), htmlFileUri)

	return nil
}
