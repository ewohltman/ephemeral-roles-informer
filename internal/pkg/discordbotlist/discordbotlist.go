// Package discordbotlist provides an implementation for updating Discord Bot List (https://top.gg).
package discordbotlist

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/DiscordBotList/go-dbl"

	"github.com/ewohltman/dbl-updater/internal/pkg/datastore"
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

	mutex                    sync.Mutex
	lastShardServerCountsSum int
}

// New returns a new *Client to update Discord Bots List.
func New(dblBotID, token string, httpClient HTTPClient, datastoreProvider datastore.Provider) (*Client, error) {
	client, err := dbl.NewClient(token, dbl.HTTPClientOption(httpClient))
	if err != nil {
		return nil, err
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
		return fmt.Errorf("error getting shard server counts from datastore provider: %w", err)
	}

	sum := 0

	for _, shardServerCount := range shardServerCounts {
		sum += shardServerCount
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if sum > client.lastShardServerCountsSum {
		client.lastShardServerCountsSum = sum

		err = client.dblClient.PostBotStats(
			client.dblBotID,
			&dbl.BotStatsPayload{
				Shards: shardServerCounts,
			},
		)
		if err != nil {
			return fmt.Errorf("error sending bot stats: %w", err)
		}

		log.Printf("Updated Discord Bot List: %d", client.lastShardServerCountsSum)
	}

	return nil
}
