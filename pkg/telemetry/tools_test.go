package telemetry

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestIsTelemetryConsent(t *testing.T) {
	t.Run("when telemetry consent is not set at all should return true", func(t *testing.T) {
		viper.Set(config.KeyConsentTelemetry.ToString(), nil)
		assert.False(t, isTelemetryConsent())
	})

	t.Run("when telemetry consent is set to false should return false", func(t *testing.T) {
		viper.Set(config.KeyConsentTelemetry.ToString(), false)
		assert.False(t, isTelemetryConsent())
	})

	t.Run("when telemetry consent is set to true should return true", func(t *testing.T) {
		viper.Set(config.KeyConsentTelemetry.ToString(), true)
		assert.True(t, isTelemetryConsent())
	})
}

func TestGetUserIdFilePath(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(homedir, ".redhat", "anonymousId"), GetUserIdFilePath())
}

func TestCreateNewUserIdentity(t *testing.T) {
	tempFilePath := fmt.Sprintf("%s/%s", os.TempDir(), "tempCreateNewUserTestFile")
	defer os.Remove(tempFilePath)

	_, err := uuid.Parse(createNewUserIdentity(tempFilePath))
	assert.NoError(t, err)
}

func TestGetCreateUserIdentity(t *testing.T) {
	t.Run("identity file exists with valid uuid should return the uuid", func(t *testing.T) {
		tempFilePath := fmt.Sprintf("%s/%s", os.TempDir(), "succesfulIdentityTestFile")
		defer os.Remove(tempFilePath)

		var expectedUid string
		require.NotPanics(t, func() { expectedUid = createNewUserIdentity(tempFilePath) })

		fetchedUid, err := GetCreateUserIdentity(tempFilePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedUid, fetchedUid)
	})

	t.Run("identity file exists without a valid uuid should create new uuid in file", func(t *testing.T) {
		tempFilePath := fmt.Sprintf("%s/%s", os.TempDir(), "notValidUidTestFile")
		defer os.Remove(tempFilePath)

		file, _ := os.Create(tempFilePath)
		defer file.Close()

		_, err := file.WriteString("not_valid_uuid")
		require.NoError(t, err)

		newUid, err := GetCreateUserIdentity(tempFilePath)
		assert.NoError(t, err)

		if _, err := uuid.Parse(newUid); err != nil {
			assert.Fail(t, err.Error())
		}

		loadedUid, _ := os.ReadFile(tempFilePath)
		assert.Equal(t, newUid, string(loadedUid))
	})

	t.Run("identify file doesn't exist should create new file with a valid uid", func(t *testing.T) {
		tempFilePath := fmt.Sprintf("%s/%s", os.TempDir(), "nonExistingTestFile")
		defer os.Remove(tempFilePath)

		if _, err := os.Stat(tempFilePath); !errors.Is(err, fs.ErrNotExist) {
			os.Remove(tempFilePath) // if the temp file exists - remove it
		}

		newUid, err := GetCreateUserIdentity(tempFilePath)
		assert.NoError(t, err)

		if _, err := uuid.Parse(newUid); err != nil {
			assert.Fail(t, err.Error())
		}

		loadedUid, _ := os.ReadFile(tempFilePath)
		assert.Equal(t, newUid, string(loadedUid))
	})
}

func TestUnwrapError(t *testing.T) {
	t.Run("unwrapping a wrapped error should return the underlying wrapped error", func(t *testing.T) {
		innerErr := fmt.Errorf("im the inner error")
		outerErr := fmt.Errorf("im the outer err and i got %w", innerErr)

		assert.Equal(t, innerErr, unwrapError(outerErr))
	})

	t.Run("unwrapping an unwrapped error should return the error itself", func(t *testing.T) {
		anErr := fmt.Errorf("dont mind me im a simple error")
		assert.Equal(t, anErr, unwrapError(anErr))
	})
}

func TestMaskErrorContent(t *testing.T) {
	template := "verify mask username to user %s and mask home folder to %s"
	expected := fmt.Sprintf(template, "$USERNAME", "$HOME")
	userInfo, _ := user.Current()
	sut := fmt.Errorf(template, userInfo.Username, userInfo.HomeDir)

	assert.Equal(t, expected, MaskErrorContent(sut))

}
