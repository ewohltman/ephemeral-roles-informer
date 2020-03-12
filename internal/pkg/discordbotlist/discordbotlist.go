// Package discordbotlist provides an implementation for updating Discord Bot List (https://top.gg).
package discordbotlist

import (
	"fmt"

	"github.com/DiscordBotList/go-dbl"

	"github.com/ewohltman/dbl-updater/internal/pkg/datasource"
)

// Client contains a *dbl.DBLClient and BotID for updating Discord Bots List.
type Client struct {
	client *dbl.DBLClient
	BotID  string
}

// New returns a new *API to update server counts.
func New(token string) (*Client, error) {
	client, err := dbl.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

// Update updates Discord Bot List with server counts obtained from the given
// datasource.Provider.
func (discordBotList Client) Update(provider datasource.Provider) error {
	shardServers, err := provider.GetShardsServerCount()
	if err != nil {
		return fmt.Errorf("error getting server counts from datastore: %w", err)
	}

	err = discordBotList.client.PostBotStats(
		discordBotList.BotID,
		dbl.BotStatsPayload{
			Shards: shardServers,
		},
	)
	if err != nil {
		return fmt.Errorf("error updating Discord Bot List: %w", err)
	}

	return nil
}
