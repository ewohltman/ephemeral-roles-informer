// Package prometheus provides a Prometheus implementation of a datasource.Provider.
package prometheus

// Prometheus contains fields for querying a Prometheus datasource.
type Prometheus struct {
}

// GetShardServers gets metrics from Prometheus and satisfies the
// datasource.Provider interface.
func (prom *Prometheus) GetShardServers() ([]int, error) {
	return make([]int, 0), nil
}
