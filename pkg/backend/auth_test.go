package backend

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestNewUserKey(t *testing.T) {
	t.Run("when request fails should return an error", func(t *testing.T) {
		// create a test server that returns a 500 status code
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the request is the one sent by this test
			if "some-fake-client1" == r.Header.Get("Client") {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		_, err := RequestNewUserKey("fake-host.org", "fake-3scale-token", "some-fake-client1")
		assert.Error(t, err)
	})

	t.Run("when request is successful should return the user id from the response", func(t *testing.T) {
		// create a test server that returns a 200 status code and a valid response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the request is the one sent by this test
			if "some-fake-client2" == r.Header.Get("Client") {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{"user_id": "t4r4e3w2"}`))
				return
			}
			w.WriteHeader(500) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		newUid, err := RequestNewUserKey(ts.URL, "fake-3scale-token", "some-fake-client2")
		assert.NoError(t, err)
		assert.Equal(t, "t4r4e3w2", newUid)

	})

}

func TestAssociateSnykToken(t *testing.T) {
	t.Run("when request fails should return an error", func(t *testing.T) {
		// create a test server that returns a 500 status code
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the request is the one sent by this test
			if "another-fake-client1" == r.Header.Get("Client") &&
				"ii88aa77hh" == r.Header.Get("Uuid") {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		err := AssociateSnykToken(
			ts.URL,
			"another-fake-3scale-token",
			"another-fake-client1",
			"ii88aa77hh",
			"oosssaakkk",
		)
		assert.Error(t, err)
	})

	t.Run("when request is successful should not return an error", func(t *testing.T) {
		// create a test server that returns a 200 status code
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the request is the one sent by this test
			if "another-fake-client2" == r.Header.Get("Client") &&
				"ii88aa77hh" == r.Header.Get("Uuid") {
				defer r.Body.Close()
				var body map[string]string
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				if body["user_id"] == "ii88aa77hh" &&
					body["snyk_api_token"] == "oosssaakkk" {
					w.WriteHeader(200)
					return
				}
			}
			w.WriteHeader(500) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		err := AssociateSnykToken(
			ts.URL,
			"another-fake-3scale-token",
			"another-fake-client2",
			"ii88aa77hh",
			"oosssaakkk",
		)
		assert.NoError(t, err)
	})
}

func TestBuildUserEndpointUrl(t *testing.T) {
	assert.Equal(
		t,
		"http://a-fake-host/user?user_key=aa11bb22",
		buildUserEndpointUrl("http://a-fake-host", "aa11bb22"),
	)
}
