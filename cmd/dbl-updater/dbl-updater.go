package main

import (
	"context"
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

	prometheusURL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
)

func update(dblClient *discordbotlist.Client) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	return dblClient.Update(ctx)
}

func main() {
	log.Printf("dbl-updater starting up")

	datastoreProvider, err := prometheus.NewProvider(prometheusURL)
	if err != nil {
		log.Fatalf("Error creating new Prometheus provider: %s", err)
	}

	dblClient, err := discordbotlist.New("", "", datastoreProvider)
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
			err = update(dblClient)
			if err != nil {
				log.Printf("Error updating Discord Bot List: %s", err)
			}
		case <-sigTerm:
			return
		}
	}
}
