// Package ds is the cloud datastore helper function
package ds

import (
	"context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

// Client type
type Client struct {
	*datastore.Client
	Cache Cache
}

// NewClient creates new ds client which wrap datastore client
func NewClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) {
	client, err := datastore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{Client: client}, nil
}
