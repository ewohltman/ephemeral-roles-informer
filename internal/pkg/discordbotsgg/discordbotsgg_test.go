package discordbotsgg

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/mock"
)

const contextTimeout = 10 * time.Second

func TestNewClient(t *testing.T) {
	dbggClient := NewClient(&http.Client{}, "dbggBotID", "dbggBotToken", &mock.Provider{})
	defer dbggClient.Close()

	if dbggClient == nil {
		t.Fatalf("Error: unexpected nil *Client")
	}
}

func TestClient_Update(t *testing.T) {
	transport := mock.DiscordBotsGGRoundTripper(t)
	httpClient := &http.Client{Transport: transport}

	dbggClient := NewClient(httpClient, "dbggBotID", "dbggBotToken", &mock.Provider{})
	defer dbggClient.Close()

	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	err := dbggClient.Update(ctx)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}
