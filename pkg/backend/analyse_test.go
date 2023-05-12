package backend

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyzeDependencyTree(t *testing.T) {
	t.Skip("WIP")
	t.Run("when request is successful should return a the request body html", func(t *testing.T) {
		dummyReport := []byte("<html><body><p>Fake Report</p></body></html>")
		// create a test server that returns a 200 status code and a html report
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the request is the one sent by this test
			if "some-client" == r.Header.Get("Client") &&
				"crda-user-id-aa11" == r.Header.Get("Uuid") &&
				"fake/contenttype" == r.Header.Get("Content-Type") {
				w.WriteHeader(200)
				_, _ = w.Write(dummyReport)
				return
			}
			w.WriteHeader(500) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		report, err := AnalyzeDependencyTree(
			ts.URL,
			"maven",
			"crda-user-id-aa11",
			"some-client",
			"fake/contenttype",
			[]byte("fake-content"),
			false,
		)

		assert.NoError(t, err)
		assert.Equal(t, dummyReport, *report) // consider asserting while converting to string for better errors
	})

	t.Run("when request fails should return an error", func(t *testing.T) {
		// create a test server that returns a 500 status code
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the request is the one sent by this test
			if "some-client2" == r.Header.Get("Client") &&
				"crda-user-id-aa12" == r.Header.Get("Uuid") &&
				"fake/contenttype2" == r.Header.Get("Content-Type") {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		_, err := AnalyzeDependencyTree(
			ts.URL,
			"maven",
			"crda-user-id-aa12",
			"some-client2",
			"fake/contenttype2",
			[]byte("fake-content"),
			false,
		)

		assert.Error(t, err)
	})
}
