package discordbotlist

import (
	"context"
	"testing"
	"time"

	"github.com/ewohltman/dbl-updater/internal/pkg/datastore"
)

const contextTimeout = 10 * time.Second

// Compile time check if *mockProvider does not satisfy the datastore.Provider
// interface.
var _ datastore.Provider = &mockProvider{}

type mockProvider struct{}

func (provider mockProvider) ProvideShardServerCounts(ctx context.Context) ([]int, error) {
	return []int{31, 28, 27, 23, 26, 19, 24, 23, 17, 28}, nil
}

func TestNew(t *testing.T) {
	_, err := New("", "", &mockProvider{})
	if err != nil {
		t.Fatalf("Error creating new Discord Bot List client: %s", err)
	}
}

func TestClient_Update(t *testing.T) {
	dblClient, err := New("", "", &mockProvider{})
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
