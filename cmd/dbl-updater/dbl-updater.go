package main

import (
	"log"

	"github.com/ewohltman/dbl-updater/internal/pkg/datasource/prometheus"
	"github.com/ewohltman/dbl-updater/internal/pkg/discordbotlist"
)

func main() {
	dblClient, err := discordbotlist.New("")
	if err != nil {
		log.Fatalf("Error creating Discord Bot List client: %s", err)
	}

	datasourcePrometheus := &prometheus.Prometheus{}

	err = dblClient.Update(datasourcePrometheus)
	if err != nil {
		log.Fatalf("Error updating Discord Bot List: %s", err)
	}
}
