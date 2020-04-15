// Package mock provides mock implementations of external dependnecies.
package mock

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ewohltman/go-discordbotsgg/pkg/api"
)

// Provider is a mock datastore.Provider.
type Provider struct{}

// ProvideShardServerCounts satisfies the datastore.Provider interface and
// returns mock data.
func (provider *Provider) ProvideShardServerCounts(ctx context.Context) ([]int, error) {
	return []int{31, 28, 27, 23, 26, 19, 24, 23, 17, 28}, nil
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (rt roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// DiscordBotListRoundTripper is an http.RoundTripper to mock the response of
// the Discord Bot List API.
func DiscordBotListRoundTripper(t *testing.T) http.RoundTripper {
	return roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			defer func() {
				closeErr := req.Body.Close()
				if closeErr != nil {
					t.Errorf("Error closing request body: %s", closeErr)
				}
			}()

			_, err := io.Copy(ioutil.Discard, req.Body)
			if err != nil {
				t.Errorf("Error reading request body: %s", err)
			}

			respBody := []byte("{}")

			return &http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Header:        make(http.Header),
				Request:       req,
				ContentLength: int64(len(respBody)),
				Body:          ioutil.NopCloser(bytes.NewReader(respBody)),
			}, nil
		},
	)
}

// DiscordBotsGGRoundTripper is an http.RoundTripper to mock the response of
// the discord.bots.gg API.
func DiscordBotsGGRoundTripper(t *testing.T) http.RoundTripper {
	return roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			defer func() {
				closeErr := req.Body.Close()
				if closeErr != nil {
					t.Errorf("Error closing request body: %s", closeErr)
				}
			}()

			reqBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Errorf("Error reading request body: %s", err)
			}

			reqStats := &api.StatsUpdate{}

			err = json.Unmarshal(reqBody, reqStats)
			if err != nil {
				t.Errorf("Error unmarshaling request body: %s", err)
			}

			respStats := &api.StatsResponse{
				Stats: reqStats.Stats,
			}

			respBody, err := json.Marshal(respStats)
			if err != nil {
				t.Errorf("Error marshaling response body: %s", err)
			}

			return &http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Header:        make(http.Header),
				Request:       req,
				ContentLength: int64(len(respBody)),
				Body:          ioutil.NopCloser(bytes.NewReader(respBody)),
			}, nil
		},
	)
}
