package discordbotlist

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/mock"
)

const contextTimeout = 10 * time.Second

func TestNewClient(t *testing.T) {
	_, err := NewClient(&http.Client{}, "dblBotID", "dblBotToken", &mock.Provider{})
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestClient_Update(t *testing.T) {
	transport := mock.DiscordBotListRoundTripper(t)
	httpClient := &http.Client{Transport: transport}

	dblClient, err := NewClient(httpClient, "dblBotID", "dblBotToken", &mock.Provider{})
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	err = dblClient.Update(ctx)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}
