// Package prometheus provides a Prometheus implementation of a datasource.Provider.
package prometheus

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	contextTimeout = 10 * time.Second

	query         = `ephemeral_roles_guilds_count{pod=~"ephemeral-roles-.+"}`
	queryError    = "query Prometheus error"
	queryWarnings = "query Prometheus warnings"
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
		return nil, fmt.Errorf("%s: %w", queryError, err)
	}

	if len(warnings) > 0 {
		return nil, fmt.Errorf("%s: %s", queryWarnings, strings.Join(warnings, ", "))
	}

	resultVector, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("%s: unable to type assert result vector", queryError)
	}

	shardsServerCounts := make([]int, len(resultVector))

	for i, sample := range resultVector {
		intVal, err := strconv.Atoi(sample.Value.String())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", queryError, err)
		}

		shardsServerCounts[i] = intVal
	}

	return shardsServerCounts, nil
}
