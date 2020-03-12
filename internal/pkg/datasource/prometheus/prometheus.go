// Package prometheus provides a Prometheus implementation of a datasource.Provider.
package prometheus

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

const (
	contextTimeout = 10 * time.Second
	query          = `ephemeral_roles_guilds_count{pod=~"ephemeral-roles-.+"}`
)

// Prometheus contains fields for querying a Prometheus datasource.
type Prometheus struct {
	API v1.API
}

// GetShardsServerCount gets metrics from Prometheus and satisfies the
// datasource.Provider interface.
func (prom *Prometheus) GetShardsServerCount() ([]int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), contextTimeout)
	defer ctxCancel()

	result, warnings, err := prom.API.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("query Prometheus error: %w", err)
	}

	if len(warnings) > 0 {
		return nil, fmt.Errorf("query Prometheus warnings: %s", strings.Join(warnings, ", "))
	}

	resultBytes, err := result.Type().MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshal Prometheus results error: %w", err)
	}

	fmt.Printf("Prometheus query result: %s\n", string(resultBytes))

	return make([]int, 0), nil
}
