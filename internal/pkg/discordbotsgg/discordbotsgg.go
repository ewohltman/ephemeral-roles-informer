// Package discordbotsgg provides an implementation for updating discord.bots.gg.
package discordbotsgg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ewohltman/go-discordbotsgg/pkg/api"
	"github.com/ewohltman/go-discordbotsgg/pkg/discordbotsgg"

	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/datastore"
)

// HTTPClient is an interface for HTTP client implementations.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Client contains a *dbl.DBLClient, dblBotID, and datastore.Provider for updating
// Discord Bots List.
type Client struct {
	dbggClient        *discordbotsgg.Client
	dbggBotID         string
	datastoreProvider datastore.Provider

	mutex            sync.Mutex
	lastServerCounts int
}

// NewClient returns a new *Client to update Discord Bots List.
func NewClient(httpClient HTTPClient, dbggBotID, dbggToken string, datastoreProvider datastore.Provider) *Client {
	dbggClient := discordbotsgg.NewClient(httpClient, dbggToken)

	return &Client{
		dbggClient:        dbggClient,
		dbggBotID:         dbggBotID,
		datastoreProvider: datastoreProvider,
	}
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
		client.lastServerCounts = serverCounts

		for i := range shardServerCounts {
			statsResponse, err := client.dbggClient.UpdateWithContext(
				ctx, client.dbggBotID, &api.StatsUpdate{
					Stats: &api.Stats{
						GuildCount: shardServerCounts[i],
						ShardCount: len(shardServerCounts),
					},
					ShardID: i,
				},
			)
			if err != nil {
				return fmt.Errorf("unable to update bot stats: %w", err)
			}

			log.Printf("Updated discord.bots.gg shard %d: %+v", i, statsResponse)
		}

		log.Printf("Updated discord.bots.gg: %d", client.lastServerCounts)
	}

	return nil
}
