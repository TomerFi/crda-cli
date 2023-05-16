package analyse

import (
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"os"
)

// TreeProvider is used to contract dependency trees providers
// i.e. Java-Maven, Node-JS
type TreeProvider interface {
	// Provide is used for providing a dependency tree that will be used as the backend request body content
	// it is also in charge of providing the body content type
	// should return an error if failed to create the dependency tree
	Provide(ctx context.Context, manifestPath string) ([]byte, string, error)
}

// Manifest is used as a type for binding a file and ecosystem names with a tree provider
type Manifest struct {
	Filename, Ecosystem string
	TreeProvider
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

// IsSupportedManifestPath is used to load a manifest file from the OS and verify we can support it
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
