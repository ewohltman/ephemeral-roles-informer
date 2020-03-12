package main

import (
	"log"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/ewohltman/dbl-updater/internal/pkg/datasource/prometheus"
)

const promURL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"

func main() {
	promClient, err := api.NewClient(api.Config{Address: promURL})
	if err != nil {
		log.Fatalf("Error creating new Prometheus API client: %s", err)
	}

	datasourcePrometheus := &prometheus.Prometheus{API: v1.NewAPI(promClient)}

	_, err = datasourcePrometheus.GetShardsServerCount()
	if err != nil {
		log.Fatalf("Error getting shards server count: %s", err)
	}
}
