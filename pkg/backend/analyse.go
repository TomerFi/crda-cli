package backend

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// AnalyzeDependencyTree is used to create the stack report against the backend
// will return the response body or an error
func AnalyzeDependencyTree(backendHost, ecosystem, crdaKey, cliClient, contentType string, content []byte, jsonOut bool) (*[]byte, error) {
	apiUrl := fmt.Sprintf("%s/api/v3/dependency-analysis/%s", backendHost, ecosystem)

	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	accept := "text/html"
	if jsonOut {
		accept = "application/json"
	}

	request.Header.Add("Client", cliClient)
	request.Header.Add("Uuid", crdaKey)
	request.Header.Add("Content-Type", contentType)
	request.Header.Add("Accept", accept)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("analyze dependencies request failed, %s", response.Status)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
