package config

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type configKey string

func (k configKey) ToString() string {
	return fmt.Sprint(k)
}

const (
	KeyConsentTelemetry configKey = "consent_telemetry"
	KeyBackendHost      configKey = "crda_backend_host"
)

var (
	configPath        = ".crda"
	configName        = "config"
	configType        = "yaml"
	configHostDefault = "http://crda-backend-crda.apps.sssc-cl01.appeng.rhecoeng.com"
)

var KnownConfigKeyStrings = []string{
	KeyConsentTelemetry.ToString(),
	KeyBackendHost.ToString(),
}

// Load is used for loading crda config from either
// the environment variables or from the $HOME/.crda/config.yaml
// returns error when failed loading/populating the config file
func Load(configDirectory string) error {
	utils.Logger.Debugf("loading config %s", configDirectory)
	// set config file from user home
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configDirectory)
	// this means viper will look for env vars before config file
	// i.e. viper.GetString("crda_backend_host") will first look for a CRDA_BACKEND_HIST env var
	viper.AutomaticEnv()
	// set defaults
	viper.SetDefault(KeyBackendHost.ToString(), configHostDefault)

	// load config and create a new file if one doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			utils.Logger.Debug("config file not found, creating a new one")
			// verify file not found
			configFilePath := buildConfigFilePath(configDirectory)
			// create folders
			if err := os.MkdirAll(configDirectory, os.ModePerm); err != nil {
				return fmt.Errorf("error creating config path: %w", err)
			}
			// create file
			if _, err := os.Create(configFilePath); err != nil {
				return fmt.Errorf("error creating config file: %w", err)
			}
			utils.Logger.Debugf("new config file created: %s", configFilePath)
		} else {
			return fmt.Errorf("error loading config file: %w", err)
		}
	}

	// TODO figure out another way to do this with constantly re-writing the file
	// write updated config to file
	utils.Logger.Debug("writing new config to file")
	if err := viper.WriteConfig(); err != nil {
		utils.Logger.Debugf("failed to write config to file, %e", err)
	}

	return nil
}

// GetConfigDirectoryPath is used to get the configuration folder for crda
// $HOME/.crda
func GetConfigDirectoryPath() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, configPath)
}

// buildConfigFilePath is used to join into a path the config file name and type with the given folder
func buildConfigFilePath(configDirectory string) string {
	return filepath.Join(configDirectory, fmt.Sprintf("%s.%s", configName, configType))
}
