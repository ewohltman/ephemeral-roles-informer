// Package prometheus provides a Prometheus implementation of a datastore.Provider.
package prometheus

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/ewohltman/ephemeral-roles-informer/internal/pkg/datastore"
)

const (
	query = `ephemeral_roles_guilds_count`

	queryError    = "error querying Prometheus"
	queryWarnings = "warning querying Prometheus"
)

// Compile time check if *Provider does not satisfy the datastore.Provider
// interface.
var _ datastore.Provider = &Provider{}

// Provider provides methods for querying a Prometheus server.
type Provider struct {
	API v1.API
}

// NewProvider returns a new *Provider for querying a Prometheus server.
func NewProvider(prometheusURL string) (*Provider, error) {
	client, err := api.NewClient(api.Config{Address: prometheusURL})
	if err != nil {
		return nil, fmt.Errorf("unable to create new Prometheus provider: %w", err)
	}

	return &Provider{API: v1.NewAPI(client)}, nil
}

// ProvideShardServerCounts queries metrics from a Prometheus server and
// satisfies the datastore.Provider interface.
func (prom *Provider) ProvideShardServerCounts(ctx context.Context) ([]int, error) {
	result, warnings, err := prom.API.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", queryError, err)
	}

	if len(warnings) > 0 {
		return nil, fmt.Errorf("%s: %s", queryWarnings, strings.Join(warnings, ", "))
	}

	resultVector, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("%s: unable to type assert result vector", queryError)
	}

	return prom.convertResults(resultVector)
}

func (prom *Provider) convertResults(resultVector model.Vector) ([]int, error) {
	shardServerCounts := make([]int, len(resultVector))

	for i, sample := range resultVector {
		intVal, err := strconv.Atoi(sample.Value.String())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", queryError, err)
		}

		shardServerCounts[i] = intVal
	}

	return shardServerCounts, nil
}
