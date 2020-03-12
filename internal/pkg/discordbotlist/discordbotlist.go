// Package discordbotlist provides an implementation for updating Discord Bot List (https://top.gg).
package discordbotlist

import (
	"context"
	"fmt"
	"log"

	"github.com/DiscordBotList/go-dbl"

	"github.com/ewohltman/dbl-updater/internal/pkg/datasource"
)

// Client contains a *dbl.DBLClient and botID for updating Discord Bots List.
type Client struct {
	dblClient          *dbl.DBLClient
	botID              string
	datasourceProvider datasource.Provider
}

// New returns a new *API to update server counts.
func New(botID, token string, datasourceProvider datasource.Provider) (*Client, error) {
	client, err := dbl.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &Client{
		dblClient:          client,
		botID:              botID,
		datasourceProvider: datasourceProvider,
	}, nil
}

// Update updates Discord Bot List with server counts obtained from the given
// datasource.Provider.
func (client *Client) Update(ctx context.Context) error {
	shardServerCounts, err := client.datasourceProvider.ProvideShardServerCounts(ctx)
	if err != nil {
		return fmt.Errorf("error getting shard server counts from datastore provider: %w", err)
	}

	log.Printf("Shard server counts: %v", shardServerCounts)

	// nolint:gocritic // will enable this later
	/*err = client.dblClient.PostBotStats(
		client.botID,
		dbl.BotStatsPayload{
			Shards: shardServerCounts,
		},
	)
	if err != nil {
		return fmt.Errorf("error sending bot stats: %w", err)
	}*/

	return nil
}
