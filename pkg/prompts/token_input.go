package prompts

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"strings"
)

type promptParts struct {
	Label, Desc string
}

var tokenPrompt = promptui.Prompt{
	Label: &promptParts{
		Label: "Snyk Token",
		Desc:  "[press Enter to continue]",
	},
	Templates: &promptui.PromptTemplates{
		Valid:   "{{ .Label | green | bold }} {{.Desc | faint }}: ",
		Invalid: "{{ .Label | red }} {{ .Desc | faint}}: ",
	},
	Validate:    validateInput,
	HideEntered: true,
}

const snykTokenUrl = "https://app.snyk.io/redhat/snyk-token"

// SnykTokenPrompt is used to prompt users for snyk tokens and return it
// returns error if failed
func SnykTokenPrompt() (string, error) {
	snykTokenUrl := fmt.Sprintf("To get Snyk Token, Please click %s", snykTokenUrl)
	fmt.Println(snykTokenUrl)
	if token, err := tokenPrompt.Run(); err != nil || token == "" {
		return "", fmt.Errorf("unable to read snyk token, try later")
	} else {
		return token, nil
	}
}

// validateInput is used to validate the token prompt input based on snyk's regex
// returns error if the input didn't match the regex
// empty tokens are valid
func validateInput(input string) error {
	token := strings.TrimSpace(input)
	if token == "" {
		utils.Logger.Debug("empty token provided")
		return nil
	}
	if utils.MatchSnykRegex(token) {
		return nil
	}
	return fmt.Errorf("invalid snyk token")
}
