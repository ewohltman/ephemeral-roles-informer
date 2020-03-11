// Package prometheus provides a Prometheus implementation of a datasource.Provider.
package prometheus

// Prometheus contains fields for querying a Prometheus datasource.
type Prometheus struct {
}

// GetMetrics gets metrics from Prometheus and satisfies the
// datasource.Provider interface.
func (prom *Prometheus) GetMetrics() (int, error) {
	return 0, nil
}
