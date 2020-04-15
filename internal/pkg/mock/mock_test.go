package mock

import (
	"bytes"
	"context"
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
	// TODO
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

func testRoundTripper(rt http.RoundTripper) error {
	client := &http.Client{Transport: rt}

	req, err := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewReader([]byte("{}")))
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
