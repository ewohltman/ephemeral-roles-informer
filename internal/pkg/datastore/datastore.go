// Package datastore provides an abstraction point for backend data providers.
package datastore

import "context"

// Provider is an interface for abstracting backend data providers.
type Provider interface {
	ProvideShardServerCounts(context.Context) ([]int, error)
}
