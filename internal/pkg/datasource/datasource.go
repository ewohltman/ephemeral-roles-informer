// Package datasource provides an abstraction point for backend data providers.
package datasource

// Provider is an interface for abstracting backend data providers.
type Provider interface {
	GetMetrics() (int, error)
}
