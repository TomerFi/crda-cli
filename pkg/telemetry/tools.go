package telemetry

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetCreateUserIdentity is used for loading or creating a new user id for telemetry reporting
// will return error if failed loading/creating the uuid for telemetry
func GetCreateUserIdentity(userIdFile string) (string, error) {
	// create path if one doesn't exist
	if err := os.MkdirAll(filepath.Dir(userIdFile), os.ModePerm); err != nil {
		return "", err
	}
	// check if file exists and create a new one if needed
	_, err := os.Stat(userIdFile)
	if errors.Is(err, fs.ErrNotExist) {
		return createNewUserIdentity(userIdFile), nil
	}
	// if statistics error is other than not found
	if err != nil {
		return "", err
	}
	// read existing file
	var id []byte
	if id, err = os.ReadFile(userIdFile); err != nil {
		return "", err
	}
	// verify file content is uuid-parsable
	var uid uuid.UUID
	if uid, err = uuid.Parse(strings.TrimSpace(string(id))); err != nil {
		return createNewUserIdentity(userIdFile), nil
	}
	// return the parsed uid
	return uid.String(), nil
}

// createNewUserIdentity is used to create a new user id file containing a random uid
// will return the uuid regardless to file write process success
func createNewUserIdentity(userIdentityFilePath string) string {
	// create a random uuid, write it to the file, and return it
	newUid := uuid.NewString()
	if err := os.WriteFile(userIdentityFilePath, []byte(newUid), 0600); err != nil {
		utils.Logger.Debugf("failed writing new uid to file, %e", err)
	}
	return newUid
}

// GetUserIdFilePath is used to construct the path for the telemetry user id
// the user id is shared between all application using red hat telemetry
// i.e. vscode rh extensions, intellij rh plugins, etc.
// $HOME/.redhat/anonymousId
func GetUserIdFilePath() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, ".redhat", "anonymousId")
}

// isTelemetryConsent will return true if the telemetry consent config value is set to true
// if value not yet set, will return false (auth and analyse commands set telemetry consent)
func isTelemetryConsent() bool {
	return viper.IsSet(config.KeyConsentTelemetry.ToString()) &&
		viper.GetBool(config.KeyConsentTelemetry.ToString())
}

func MaskErrorContent(err error) string {
	user, usrErr := user.Current()
	if usrErr != nil {
		return usrErr.Error()
	}

	sanitized := unwrapError(err).Error()
	sanitized = strings.ReplaceAll(sanitized, user.HomeDir, "$HOME")
	sanitized = strings.ReplaceAll(sanitized, user.Username, "$USERNAME")

	return sanitized
}

// unwrapError is used to extract the wrapped error if one exists
// will return either the internal error if one exists or the original one
func unwrapError(err error) error {
	wrapped := errors.Unwrap(err)
	if wrapped != nil {
		return wrapped
	}
	return err
}
