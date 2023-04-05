package auth

import (
	"context"
	"encoding/json"
	"github.com/rhecosystemappeng/crda-cli/pkg/config"
	"github.com/rhecosystemappeng/crda-cli/pkg/telemetry"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	utils.ConfigureLogging(false)
}

func TestAuthenticateUser(t *testing.T) {
	validSnykToken := "a012345B-123C-432d-B123-80123456789B" // if token not valid, input is requested

	t.Run("authenticate user with an existing crda key (uuid)", func(t *testing.T) {
		// create a test server that return a 200 status code and a valid response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the request is the one sent by this test
			defer r.Body.Close()
			var body map[string]string
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			if body["user_id"] == "a-fake-crda-key1" &&
				body["snyk_api_token"] == validSnykToken {
				w.WriteHeader(200)
				return
			}
			w.WriteHeader(500) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		ctx := telemetry.GetContext(context.Background())
		telemetry.SetProperty(ctx, telemetry.KeyClient, "lets-say-im-intellij")

		viper.Set(config.KeyConsentTelemetry.ToString(), false) // if this is nil, input is requested
		viper.Set(config.KeyOldHost.ToString(), ts.URL)         // route requests to the test server
		viper.Set(config.KeyOld3ScaleToken.ToString(), "aaoo9988aa77jj")

		// if we have an uuid, a new one won't be requested from the server
		viper.Set(config.KeyCrdaKey.ToString(), "a-fake-crda-key1")

		assert.NoError(t, AuthenticateUser(ctx, validSnykToken))
	})

	t.Run("authenticate user with no crda key (uuid) should request one from the server", func(t *testing.T) {
		// create a test server that return a 200 status code and a valid response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// verify the requests are the ones sent by this test
			if "POST" == r.Method {
				// POST requests are sent when a new crda uuid is required to be issued by the server
				if "lets-say-im-intellij" == r.Header.Get("Client") {
					w.WriteHeader(200)
					_, _ = w.Write([]byte(`{"user_id": "a-fake-crda-key2"}`))
					return
				}
			} else if "PUT" == r.Method {
				// PUT requests are sent when a snyk token is required be associated with a crda uuid
				defer r.Body.Close()
				var body map[string]string
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				if body["user_id"] == "a-fake-crda-key2" &&
					body["snyk_api_token"] == validSnykToken {
					w.WriteHeader(200)
					return
				}
			}
			w.WriteHeader(500) // this shouldn't happen and will fail the test (intentionally)
		}))
		defer ts.Close()

		ctx := telemetry.GetContext(context.Background())
		telemetry.SetProperty(ctx, telemetry.KeyClient, "lets-say-im-intellij")

		viper.Set(config.KeyConsentTelemetry.ToString(), false) // if this is nil, input is requested
		viper.Set(config.KeyOldHost.ToString(), ts.URL)         // route requests to the test server
		viper.Set(config.KeyOld3ScaleToken.ToString(), "aaoo9988aa77jj")

		// if we have an uuid, a new one won't be requested from the server
		viper.Set(config.KeyCrdaKey.ToString(), nil)

		assert.NoError(t, AuthenticateUser(ctx, validSnykToken))
	})
}
