package analyse

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestGetManifest(t *testing.T) {
	for _, manifest := range SupportedManifests {
		t.Run(fmt.Sprintf("succesfully get manifest for %s", manifest.Filename), func(t *testing.T) {
			fetched, err := GetManifest(manifest.Filename)
			assert.NoError(t, err)
			assert.Equal(t, manifest.Filename, fetched.Filename)
			assert.Equal(t, manifest.Ecosystem, fetched.Ecosystem)
		})
	}

	t.Run("failed to get manifest for unknown package file", func(t *testing.T) {
		_, err := GetManifest("no-a-real.file")
		assert.Error(t, err)
	})
}

func TestInitialization(t *testing.T) {
	t.Run("verify all manifest types are includes in the string slice", func(t *testing.T) {
		for _, manifest := range SupportedManifests {
			t.Run(fmt.Sprintf("verifying %s exists", manifest.Filename), func(t *testing.T) {
				assert.True(t, slices.Contains(SupportedManifestsFilenames, manifest.Filename))
			})
		}
	})
}

func TestIsSupportedManifestPath(t *testing.T) {
	for _, manifestName := range SupportedManifestsFilenames {
		t.Run(fmt.Sprintf("verifing %s is supported should not return an error", manifestName), func(t *testing.T) {
			tempManifest := filepath.Join(os.TempDir(), manifestName)
			defer os.Remove(tempManifest)

			if _, err := os.Stat(tempManifest); !errors.Is(err, fs.ErrNotExist) {
				os.Remove(tempManifest) // if the temp file exists - remove it
			}

			_, err := os.Create(tempManifest)
			require.NoError(t, err)

			assert.NoError(t, IsSupportedManifestPath(tempManifest))
		})
	}

	t.Run("verifying a directory should return an error", func(t *testing.T) {
		assert.Error(t, IsSupportedManifestPath(os.TempDir()))
	})
}
