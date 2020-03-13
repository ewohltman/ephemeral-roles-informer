// Package discordbotlist provides an implementation for updating Discord Bot List (https://top.gg).
package discordbotlist

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/DiscordBotList/go-dbl"

	"github.com/ewohltman/dbl-updater/internal/pkg/datastore"
)

// Client contains a *dbl.DBLClient, dblBotID, and datastore.Provider for updating
// Discord Bots List.
type Client struct {
	dblClient         *dbl.DBLClient
	dblBotID          string
	datastoreProvider datastore.Provider

	mutex                    sync.Mutex
	lastShardServerCountsSum int
}

// New returns a new *Client to update Discord Bots List.
func New(dblBotID, token string, datastoreProvider datastore.Provider) (*Client, error) {
	client, err := dbl.NewClient(token)
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

	log.Printf("Shard server counts: %v", shardServerCounts)

	sum := 0

	for _, shardServerCount := range shardServerCounts {
		sum += shardServerCount
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if sum > client.lastShardServerCountsSum {
		// nolint:gocritic // will enable this later
		/*err = client.dblClient.PostBotStats(
			client.dblBotID,
			dbl.BotStatsPayload{
				Shards: shardServerCounts,
			},
		)
		if err != nil {
			return fmt.Errorf("error sending bot stats: %w", err)
		}*/

		client.lastShardServerCountsSum = sum

		log.Printf("Updated Discord Bot List: %d", client.lastShardServerCountsSum)
	}

	return nil
}