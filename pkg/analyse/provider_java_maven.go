package analyse

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
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

type JavaMavenTreeProvider struct{}

func (a *JavaMavenTreeProvider) Provide(ctx context.Context, manifestPath string) ([]byte, string, error) {
	mvn, err := exec.LookPath("mvn")
	if err != nil {
		return nil, "", err
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
	cleanExec := exec.CommandContext(ctx, mvn, "-q", "clean", "-f", manifestPath)
	treeExec := exec.CommandContext(ctx, mvn, treeCommand...)

	if err := cleanExec.Run(); err != nil {
		return nil, "", err
	}

	if err := treeExec.Run(); err != nil {
		return nil, "", err
	}

	graph, err := os.ReadFile(tmpDesTree)
	if err != nil {
		return nil, "", err
	}

	return graph, "text/vnd.graphviz", nil
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
