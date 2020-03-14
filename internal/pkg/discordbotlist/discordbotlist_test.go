package discordbotlist

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ewohltman/dbl-updater/internal/pkg/datastore"
)

const contextTimeout = 10 * time.Second

// RoundTripperFunc allows functions to satisfy the http.RoundTripper
// interface.
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

// RoundTrip implements the http.RoundTripper interface.
func (rt RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// Compile time check if *mockProvider does not satisfy the datastore.Provider
// interface.
var _ datastore.Provider = &mockProvider{}

type mockProvider struct{}

func (provider mockProvider) ProvideShardServerCounts(ctx context.Context) ([]int, error) {
	return []int{31, 28, 27, 23, 26, 19, 24, 23, 17, 28}, nil
}

func TestNew(t *testing.T) {
	_, err := New("botID", "token", &http.Client{}, &mockProvider{})
	if err != nil {
		t.Fatalf("Error creating new Discord Bot List client: %s", err)
	}
}

func TestClient_Update(t *testing.T) {
	transport := mockDiscordBotListAPI(t)
	httpClient := &http.Client{Transport: transport}

	dblClient, err := New("botID", "token", httpClient, &mockProvider{})
	if err != nil {
		t.Fatalf("Error creating new Discord Bot List client: %s", err)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	err = dblClient.Update(ctx)
	if err != nil {
		t.Fatalf("Error updating Discord Bot List: %s", err)
	}
}

func mockDiscordBotListAPI(t *testing.T) http.RoundTripper {
	return RoundTripperFunc(
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
