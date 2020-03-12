// Package prometheus provides a Prometheus implementation of a datasource.Provider.
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

	"github.com/ewohltman/dbl-updater/internal/pkg/datasource"
)

const (
	prometheusURL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
	query         = `ephemeral_roles_guilds_count{pod=~"ephemeral-roles-.+"}`

	queryError    = "query Prometheus error"
	queryWarnings = "query Prometheus warnings"
)

// Compile time error if *Prometheus does not satisfy the datasource.Provider
// interface.
var _ datasource.Provider = &Prometheus{}

// Prometheus provides methods for querying a Prometheus datasource and
// satisfies the datasource.Provider interface.
type Prometheus struct {
	API v1.API
}

// New returns a new *Prometheus instance for querying Prometheus metrics.
func New() (*Prometheus, error) {
	client, err := api.NewClient(api.Config{Address: prometheusURL})
	if err != nil {
		return nil, err
	}

	return &Prometheus{API: v1.NewAPI(client)}, nil
}

// ProvideShardServerCounts gets metrics from Prometheus and satisfies the
// datasource.Provider interface.
func (prom *Prometheus) ProvideShardServerCounts(ctx context.Context) ([]int, error) {
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
