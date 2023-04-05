package utils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"os"
	"strings"
	"testing"
)

func TestMatchSnykRegex(t *testing.T) {
	t.Run("empty string as token should return false", func(t *testing.T) {
		assert.False(t, MatchSnykRegex(""))
	})

	t.Run("legit token should return true", func(t *testing.T) {
		assert.True(t, MatchSnykRegex("a012345B-123C-432d-B123-80123456789B"))
	})
}

func TestReportBodyToTempHtmlFile(t *testing.T) {
	// save html content to html file
	dummyReport := []byte("<html><body><p>Fake Report</p></body></html>")
	uri, err := SaveReportToTempHtmlFile(dummyReport, "testecosystem")
	require.NoError(t, err)

	// read html file
	content, err := os.ReadFile(strings.Replace(uri, "file://", "", 1))
	require.NoError(t, err)

	// parse html file
	_, err = html.Parse(bytes.NewReader(content))
	require.NoError(t, err)

	// verify html content
	assert.Equal(t, dummyReport, content)
}
