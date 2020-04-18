// Package discordbotlist provides an implementation for updating Discord Bot List (https://top.gg).
package discordbotlist

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/DiscordBotList/go-dbl"

	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/datastore"
)

// HTTPClient is an interface for HTTP client implementations.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Client contains a *dbl.DBLClient, dblBotID, and datastore.Provider for updating
// Discord Bots List.
type Client struct {
	dblClient         *dbl.Client
	dblBotID          string
	datastoreProvider datastore.Provider

	mutex            sync.Mutex
	lastServerCounts int
}

// NewClient returns a new *Client to update Discord Bots List.
func NewClient(httpClient HTTPClient, dblBotID, dblBotToken string, datastoreProvider datastore.Provider) (*Client, error) {
	client, err := dbl.NewClient(dblBotToken, dbl.HTTPClientOption(httpClient))
	if err != nil {
		return nil, fmt.Errorf("unable to create new discord bot list client: %w", err)
	}

	return &Client{
		dblClient:         client,
		dblBotID:          dblBotID,
		datastoreProvider: datastoreProvider,
	}, nil
}

// Update updates Discord Bot List with server counts obtained from a
// datastore.Provider.
func (client *Client) Update(ctx context.Context) error {
	shardServerCounts, err := client.datastoreProvider.ProvideShardServerCounts(ctx)
	if err != nil {
		return fmt.Errorf("unable to get shard server counts from datastore provider: %w", err)
	}

	serverCounts := 0

	for _, shardServerCount := range shardServerCounts {
		serverCounts += shardServerCount
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if serverCounts > client.lastServerCounts {
		err = client.dblClient.PostBotStats(
			client.dblBotID,
			&dbl.BotStatsPayload{
				Shards: shardServerCounts,
			},
		)
		if err != nil {
			return fmt.Errorf("unable to update bot stats: %w", err)
		}

		client.lastServerCounts = serverCounts

		log.Printf("Updated Discord Bot List: %d", client.lastServerCounts)
	}

	return nil
}
