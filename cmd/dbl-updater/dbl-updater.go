package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ewohltman/dbl-updater/internal/pkg/datastore/prometheus"
	"github.com/ewohltman/dbl-updater/internal/pkg/discordbotlist"
)

const (
	updateInterval = 30 * time.Second
	contextTimeout = 10 * time.Second

	envBotID = "DBL_BOT_ID"
	envToken = "DBL_BOT_TOKEN" // nolint:gosec // Credential from environment variable, not hardcoded

	prometheusURL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"

	errEnvNotFound = "%s not defined"
)

func environmentLookup() (dblBotID, token string, err error) {
	var found bool

	dblBotID, found = os.LookupEnv(envBotID)
	if !found {
		return "", "", fmt.Errorf(errEnvNotFound, envBotID)
	}

	token, found = os.LookupEnv(envToken)
	if !found {
		return "", "", fmt.Errorf(errEnvNotFound, envToken)
	}

	return
}

func updateDiscordBotList(dblClient *discordbotlist.Client) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	return dblClient.Update(ctx)
}

func main() {
	log.Printf("dbl-updater starting up")

	dblBotID, token, err := environmentLookup()
	if err != nil {
		log.Fatalf("Error looking up environment variables: %s", err)
	}

	datastoreProvider, err := prometheus.NewProvider(prometheusURL)
	if err != nil {
		log.Fatalf("Error creating new Prometheus provider: %s", err)
	}

	dblClient, err := discordbotlist.New(dblBotID, token, datastoreProvider)
	if err != nil {
		log.Fatalf("Error creating new Discord Bot List client: %s", err)
	}

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = updateDiscordBotList(dblClient)
			if err != nil {
				log.Printf("Error updating Discord Bot List: %s", err)
			}
		case <-sigTerm:
			return
		}
	}
}
