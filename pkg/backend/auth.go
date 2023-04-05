package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestNewUserKey is used to requesting a new uuid from the backend server
func RequestNewUserKey(host, tsToken, cliClient string) (string, error) {
	apiUrl := buildUserEndpointUrl(host, tsToken)

	request, err := http.NewRequest(http.MethodPost, apiUrl, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Client", cliClient)
	request.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("new user key request failed, %s", response.Status)
	}

	type newUserKeyResponse struct {
		Key string `json:"user_id,omitempty"`
	}

	var parsed newUserKeyResponse
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return "", err
	}

	return parsed.Key, nil
}

// AssociateSnykToken is used for requesting the backed server to associate a snyk token with an uuid
func AssociateSnykToken(host, tsToken, cliClient, crdaKey, snykToken string) error {
	apiUrl := buildUserEndpointUrl(host, tsToken)

	type associateSnykTokenPayload struct {
		CrdaKey   string `json:"user_id"`
		SnykToken string `json:"snyk_api_token"`
	}

	body, err := json.Marshal(associateSnykTokenPayload{crdaKey, snykToken})
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Add("Client", cliClient)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Uuid", crdaKey)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("associate snyk token request failed, %s", response.Status)
	}

	return nil
}

func buildUserEndpointUrl(host, tsToken string) string {
	return fmt.Sprintf("%s/user?user_key=%s", host, tsToken)
}
