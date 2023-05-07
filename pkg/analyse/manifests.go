package analyse

import (
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"os"
)

type TreeProvider interface {
	Provide(ctx context.Context, manifestPath string) ([]byte, string, error)
}

type Manifest struct {
	Filename, Ecosystem string
	TreeProvider
}

type Provider string

const (
	ProviderSnyk Provider = "snyk"
)

func (p Provider) ToString() string {
	return fmt.Sprint(p)
}

var (
	JavaMaven = Manifest{"pom.xml", "maven", &JavaMavenTreeProvider{}}
	PythonPip = Manifest{"requirements.txt", "maven", nil}
	NodeJS    = Manifest{"package.json", "npm", nil}
	GoModule  = Manifest{"go.mod", "go", nil}
)

var SupportedManifests = []Manifest{JavaMaven, PythonPip, NodeJS, GoModule}
var SupportedManifestsFilenames []string

func init() {
	// create a string slice for coordinating the supported package file names
	for _, m := range SupportedManifests {
		SupportedManifestsFilenames = append(SupportedManifestsFilenames, m.Filename)
	}
}

// GetManifest returns the Manifest type for a string
// returns error then used with an unknown manifest file
func GetManifest(fileName string) (*Manifest, error) {
	for _, m := range SupportedManifests {
		if m.Filename == fileName {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("manifest %s not supported", fileName)
}

func IsSupportedManifestPath(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("please use path to manifest file as input")
	}
	if !slices.Contains(SupportedManifestsFilenames, fileInfo.Name()) {
		return fmt.Errorf("manifest %s is not supported", fileInfo.Name())
	}
	return nil
}
