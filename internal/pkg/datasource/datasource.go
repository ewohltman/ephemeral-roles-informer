// Package datasource provides an abstraction point for backend data providers.
package datasource

import "context"

// Provider is an interface for abstracting backend data providers.
type Provider interface {
	ProvideShardServerCounts(context.Context) ([]int, error)
}
