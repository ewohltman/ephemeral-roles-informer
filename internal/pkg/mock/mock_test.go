package mock

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestProvider_ProvideShardServerCounts(t *testing.T) {
	provider := &Provider{}

	_, err := provider.ProvideShardServerCounts(context.Background())
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestRoundTripperFunc_RoundTrip(t *testing.T) {
	err := testRoundTripper(DiscordBotListRoundTripper(t))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestDiscordBotListRoundTripper(t *testing.T) {
	err := testRoundTripper(DiscordBotListRoundTripper(t))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestDiscordBotsGGRoundTripper(t *testing.T) {
	err := testRoundTripper(DiscordBotsGGRoundTripper(t))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func testRoundTripper(rt http.RoundTripper) (err error) {
	client := &http.Client{Transport: rt}

	req, err := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewReader([]byte("{}")))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%s: %w", closeErr, err)
				return
			}

			err = closeErr
		}
	}()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
