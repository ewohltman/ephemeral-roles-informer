package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/datastore/prometheus"
	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/discordbotlist"
	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/discordbotsgg"
)

const (
	connectionTimeout = 10 * time.Second
	contextTimeout    = 20 * time.Second
	updateInterval    = 30 * time.Second

	envDBLBotID     = "DBL_BOT_ID"
	envDBLBotToken  = "DBL_BOT_TOKEN" // nolint:gosec // Credential from environment variable, not hardcoded
	envDBGGBotID    = "DBGG_BOT_ID"
	envDBGGBotToken = "DBGG_BOT_TOKEN" // nolint:gosec // Credential from environment variable, not hardcoded

	prometheusURL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"

	errEnvNotFound = "unable to lookup %s: not defined in environment variables"
)

func dblEnvironmentLookup() (dblBotID, dblBotToken string, err error) {
	var found bool

	dblBotID, found = os.LookupEnv(envDBLBotID)
	if !found {
		return "", "", fmt.Errorf(errEnvNotFound, envDBLBotID)
	}

	dblBotToken, found = os.LookupEnv(envDBLBotToken)
	if !found {
		return "", "", fmt.Errorf(errEnvNotFound, envDBLBotToken)
	}

	return
}

func dbggEnvironmentLookup() (dbggBotID, dbggBotToken string, err error) {
	var found bool

	dbggBotID, found = os.LookupEnv(envDBGGBotID)
	if !found {
		return "", "", fmt.Errorf(errEnvNotFound, envDBGGBotID)
	}

	dbggBotToken, found = os.LookupEnv(envDBGGBotToken)
	if !found {
		return "", "", fmt.Errorf(errEnvNotFound, envDBGGBotToken)
	}

	return
}

func updateDiscordBotList(dblClient *discordbotlist.Client) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	return dblClient.Update(ctx)
}

func updateDiscordBotsGG(dbggClient *discordbotsgg.Client) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	return dbggClient.Update(ctx)
}

func main() {
	log.Printf("ephemeral-roles-informer starting up")

	datastoreProvider, err := prometheus.NewProvider(prometheusURL)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	httpClient := &http.Client{Timeout: connectionTimeout}

	dblBotID, dblBotToken, err := dblEnvironmentLookup()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	dblClient, err := discordbotlist.NewClient(httpClient, dblBotID, dblBotToken, datastoreProvider)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	dbggBotID, dbggBotToken, err := dbggEnvironmentLookup()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	dbggClient := discordbotsgg.NewClient(httpClient, dbggBotID, dbggBotToken, datastoreProvider)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				err = updateDiscordBotList(dblClient)
				if err != nil {
					log.Printf("Error updating Discord Bot List: %s", err)
				}
			}()

			go func() {
				err = updateDiscordBotsGG(dbggClient)
				if err != nil {
					log.Printf("Error updating discord.bots.gg: %s", err)
				}
			}()
		case <-sigTerm:
			return
		}
	}
}
