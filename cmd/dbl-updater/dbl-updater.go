package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/ewohltman/dbl-updater/internal/pkg/datasource/prometheus"
)

const promURL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"

func main() {
	log.Printf("dbl-updater starting up")

	promClient, err := api.NewClient(api.Config{Address: promURL})
	if err != nil {
		log.Fatalf("Error creating new Prometheus API client: %s", err)
	}

	datasourcePrometheus := &prometheus.Prometheus{API: v1.NewAPI(promClient)}

	shardsServerCounts, err := datasourcePrometheus.GetShardsServerCount()
	if err != nil {
		log.Fatalf("Error getting shards server count: %s", err)
	}

	log.Printf("Shard server counts: %v", shardsServerCounts)

	sigTerm := make(chan os.Signal, 1)

	signal.Notify(sigTerm, syscall.SIGTERM)

	<-sigTerm
}
