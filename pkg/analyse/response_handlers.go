package analyse

import (
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/backend"
	"io"
	"mime/multipart"
)

func handleJsonResponse(body io.ReadCloser) error {
	if reports, err := backend.ParseJsonResponse(body); err != nil {
		return err
	} else {
		vulSummary, err := processVulnerabilities(reports)
		if err != nil {
			return err
		}
		return printJson(vulSummary)
	}
}

func handleMixedResponse(body io.ReadCloser, params map[string]string, ecosystem string) error {
	var vulSummary VulnerabilitiesSummary
	var reportUri string

	multipartReader := multipart.NewReader(body, params["boundary"])
	for part, err := multipartReader.NextPart(); err == nil; part, err = multipartReader.NextPart() {
		partType := part.Header.Get("Content-Type")
		switch partType {
		case "application/json":
			reports, err := backend.ParseJsonResponse(part)
			if err != nil {
				return err
			}
			if vulSummary, err = processVulnerabilities(reports); err != nil {
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

	return printSummary(vulSummary, reportUri)
}
