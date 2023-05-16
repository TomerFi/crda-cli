package analyse

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend/api"
	"io"
	"mime/multipart"
)

// handleJsonResponse will deserialize the application/json dependency analysis response body
// into an api.AnalysisReport and pass it for printing
// will return error if failed to deserialize
func handleJsonResponse(body io.ReadCloser, verboseOut bool) error {
	report, err := backend.ParseJsonResponse(body)
	if err != nil {
		return err
	}
	return printJson(report, verboseOut)
}

// handleJsonResponse will deserialize the multipart/mixed dependency analysis response body
// into an api.AnalysisReport and pass it for printing
// multipart/mixed is constructed from an application/json and a text/html content parts
// will return error if failed to deserialize
func handleMixedResponse(body io.ReadCloser, params map[string]string, ecosystem string, verboseOut bool) error {
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

	if verboseOut {
		printVerboseSummary(report, reportUri)
	} else {
		printSummary(report, reportUri)
	}
	return nil
}
