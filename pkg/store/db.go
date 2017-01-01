package store

import (
	"context"
	"crypto/rsa"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauthjwt "golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

// DB type
type DB struct {
	client    *datastore.Client
	publicKey *rsa.PublicKey
}

// NewDB create new database instance
func NewDB(options ...Option) *DB {
	opts := &Options{}

	for _, setter := range options {
		setter(opts)
	}

	ctx := context.Background()

	var ts oauth2.TokenSource
	var err error

	if opts.JSONKey != nil {
		var conf *oauthjwt.Config
		conf, err = google.JWTConfigFromJSON(
			opts.JSONKey,
			datastore.ScopeDatastore,
			storage.ScopeReadWrite,
		)
		if err != nil {
			panic(err)
		}
		ts = conf.TokenSource(ctx)
	} else {
		ts, err = google.DefaultTokenSource(
			ctx,
			datastore.ScopeDatastore,
			storage.ScopeReadWrite,
		)
		if err != nil {
			panic(err)
		}
	}

	db := DB{}

	db.client, err = datastore.NewClient(ctx, opts.ProjectID, option.WithTokenSource(ts))
	if err != nil {
		panic(err)
	}

	return &db
}
