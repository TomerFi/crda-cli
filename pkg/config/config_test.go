package config

import (
	"errors"
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"regexp"
	"testing"
)

func init() {
	utils.ConfigureLogging(false)
}

func TestLoad(t *testing.T) {
	t.Run("loading with no config file should create a new one with default values", func(t *testing.T) {
		tmpConfigFolder := fmt.Sprintf("%s/crdaNewConfigTst", os.TempDir())
		defer os.RemoveAll(tmpConfigFolder)

		if _, err := os.Stat(tmpConfigFolder); !errors.Is(err, fs.ErrNotExist) {
			os.RemoveAll(tmpConfigFolder) // if the temp folder exists - remove it
		}

		require.NoError(t, Load(tmpConfigFolder)) // load will create the config file with the default values

		viper.Reset() // reset viper to forget current configuration
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
		viper.AddConfigPath(tmpConfigFolder)

		require.NoError(t, viper.ReadInConfig()) // load configuration from the saved file

		assert.Equal(t, viper.GetString(KeyBackendHost.ToString()), configHostDefault)
	})

	t.Run("loading with an existing config file should not load the default values", func(t *testing.T) {
		tmpConfigFolder := fmt.Sprintf("%s/crdaExistingConfigTst", os.TempDir())
		defer os.RemoveAll(tmpConfigFolder)

		// create config file if it doesn't exist
		if _, err := os.Stat(tmpConfigFolder); errors.Is(err, fs.ErrNotExist) {
			_ = os.MkdirAll(tmpConfigFolder, os.ModePerm)
			_, _ = os.Create(buildConfigFilePath(tmpConfigFolder))
		}

		// prepare config data to write to the file
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
		viper.AddConfigPath(tmpConfigFolder)
		viper.Set(KeyBackendHost.ToString(), "this-is-fake-host")

		// write the config to the file and reset viper so that we can start testing
		require.NoError(t, viper.WriteConfig())
		viper.Reset()

		// load the config file and verify we got the values we expect
		require.NoError(t, Load(tmpConfigFolder))
		assert.Equal(t, "this-is-fake-host", viper.GetString(KeyBackendHost.ToString()))
	})
}

func TestGetConfigDirectoryPath(t *testing.T) {
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("[a-zA-Z]/%s$", configPath)), GetConfigDirectoryPath())
}

func TestConfigKey_ToString(t *testing.T) {
	assert.Equal(t, KeyConsentTelemetry.ToString(), fmt.Sprint(KeyConsentTelemetry))
}

func TestBuildConfigFilePath(t *testing.T) {
	assert.Equal(
		t,
		fmt.Sprintf("/dummy_folder/%s.%s", configName, configType),
		buildConfigFilePath("/dummy_folder"),
	)
}
