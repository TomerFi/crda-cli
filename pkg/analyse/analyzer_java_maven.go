package analyse

import (
	"context"
	"encoding/json"
	"encoding/xml"
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
	"strings"
)

type PomProject struct {
	Dependencies PomDependencies `xml:"dependencies"`
}

type PomDependencies struct {
	Dependency []PomDependency `xml:"dependency"`
}

type PomDependency struct {
	Comment    string `xml:",comment"`
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version,omitempty"`
}

type JavaMavenAnalyzer struct{}

func (a *JavaMavenAnalyzer) Analyze(ctx context.Context, ecosystem string, manifestPath string, jsonOut, verboseOut bool) error {
	mvn, err := exec.LookPath("mvn")
	if err != nil {
		return err
	}

	tmpDesTree := filepath.Join(os.TempDir(), "tmp-deps-tree.txt")
	defer os.Remove(tmpDesTree)

	if _, err := os.Stat(tmpDesTree); !errors.Is(err, fs.ErrNotExist) {
		os.Remove(tmpDesTree) // if the temp file exists - remove it
	}

	treeCommand := []string{"-q", "dependency:tree", "-DoutputType=dot", fmt.Sprintf("-DoutputFile=%s", tmpDesTree), "-f", manifestPath}
	if ignoredList := getIgnored(manifestPath); len(ignoredList) > 0 {
		treeCommand = append(treeCommand, fmt.Sprintf("-Dexcludes=%s", strings.Join(ignoredList, ",")))
	}

	// execute commands to create a tree graph
	cleanExec := exec.Command(mvn, "-q", "clean", "-f", manifestPath)
	treeExec := exec.Command(mvn, treeCommand...)

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
		ecosystem,
		crdaKey,
		cliClient,
		"text/vnd.graphviz",
		graph,
		jsonOut,
	)
	if err != nil {
		return err
	}

	if jsonOut {
		var report []backend.DependencyAnalysisReport
		if err := json.Unmarshal(*body, &report); err != nil {
			return err
		}
		pretty, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(pretty))
		return nil
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

// getIgnored takes pom.xml path and return a list of exclusion strings for dependencies marked for ignore.
// you can add a <!-- crdaignore --> comment next to any element in the confines of <dependency>...</dependency>
// and it will be included in list returned.
// for reference, our exclusion strings looks like this "group-id:artifact-d:*:version", if no <version> element,
// will use *
func getIgnored(manifestPath string) []string {
	// load pom.xml manifest
	fileContent, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil
	}
	// deserialize the file content
	var pomProject PomProject
	if err := xml.Unmarshal(fileContent, &pomProject); err != nil {
		return nil
	}
	// populate the ignored list
	var ignoredList []string
	for _, dep := range pomProject.Dependencies.Dependency {
		if strings.Contains(dep.Comment, "crdaignore") {
			ver := "*"
			if dep.Version != "" {
				ver = dep.Version
			}
			ignoredList = append(ignoredList, fmt.Sprintf("%s:%s:*:%s", dep.GroupId, dep.ArtifactId, ver))
		}
	}

	return ignoredList
}
