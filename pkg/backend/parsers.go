package backend

import (
	"encoding/json"
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend/api"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ParseJsonResponse is used to deserialize the backend application/json type response
// It will return the unmarshalled response or an error
func ParseJsonResponse(reader io.ReadCloser) ([]api.DependencyReport, error) {
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var reports []api.DependencyReport
	if err := json.Unmarshal(body, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}

// ParseHtmlResponse is used to deserialized backend text/html type response
// It will save the unmarshalled content as a html file in os temp folder and return its uri or an error
func ParseHtmlResponse(reader io.ReadCloser, ecosystem string) (string, error) {
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	tmpFolder := filepath.Join(os.TempDir(), "crda")
	if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
		return "", err
	}

	tempFileName := fmt.Sprintf("stack-analysis-%s-%d.html", ecosystem, time.Now().Unix())
	tempFilePath := filepath.Join(tmpFolder, tempFileName)

	if err := os.WriteFile(tempFilePath, body, os.ModePerm); err != nil {
		return "", err
	}

	return fmt.Sprintf("file://%s", tempFilePath), nil
}
