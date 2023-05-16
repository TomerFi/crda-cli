package telemetry

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
)

const (
	// telemetry consent select options
	consentYes   = "Yes"
	consentNo    = "No"
	consentLater = "Later"
)

var telemetryConsent = promptui.Select{
	Label: "Would you like to contribute towards anonymous usage statistics?",
	Items: []string{consentYes, consentNo, consentLater},
}

// AskForConsent is used to prompt the user for selecting Yes|No|Later in regard to telemetry consent
// Yes or No: config will be updated accordingly, and we won't ask again
// Later: no updates to config, will ask again
//
// if consent config already recorded, will not ask at all
func AskForConsent() {
	// if consent config exists, return without asking again
	if viper.IsSet(config.KeyConsentTelemetry.ToString()) {
		return
	}
	fmt.Println("We're constantly improving this tool and would like to know more about its usage (more details at https://developers.redhat.com/article/tool-data-collection)")
	fmt.Println("Your preference can be changed manually if desired using 'crda config set consent_telemetry true|false'")

	if _, consent, err := telemetryConsent.Run(); err == nil {
		if consent == consentLater {
			fmt.Println("Ok. I will ask you again later")
		} else {
			switch consent {
			case consentYes:
				viper.Set(config.KeyConsentTelemetry.ToString(), true)
				fmt.Println("Thanks for helping us! You can disable telemetry using 'crda config set consent_telemetry false'")
			case consentNo:
				viper.Set(config.KeyConsentTelemetry.ToString(), false)
				fmt.Println("No worries, you can still enable telemetry using 'crda config set consent_telemetry true'")
			}
			if err := viper.WriteConfig(); err != nil {
				utils.Logger.Debug("failed to write configuration for telemetry consent")
			}
		}
	} else {
		utils.Logger.Debug("failed to get user consent for telemetry")
	}
}
