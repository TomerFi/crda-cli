package analyse

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend/api"
	"io"
	"mime/multipart"
)

func handleJsonResponse(body io.ReadCloser) error {
	report, err := backend.ParseJsonResponse(body)
	if err != nil {
		return err
	}
	return printJson(report)
}

func handleMixedResponse(body io.ReadCloser, params map[string]string, ecosystem string) error {
	var report *api.AnalysisReport
	var reportUri string

	multipartReader := multipart.NewReader(body, params["boundary"])
	for part, err := multipartReader.NextPart(); err == nil; part, err = multipartReader.NextPart() {
		partType := part.Header.Get("Content-Type")
		switch partType {
		case "application/json":
			report, err = backend.ParseJsonResponse(part)
			if err != nil {
				return err
			}
		case "text/html":
			if reportUri, err = backend.ParseHtmlResponse(part, ecosystem); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown response type %s", partType)
		}
	}

	printSummary(report, reportUri)
	return nil
}
