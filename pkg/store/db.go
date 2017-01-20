package store

import (
	"context"
	"crypto/rsa"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/acoshift/ds"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauthjwt "golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

// DB type
type DB struct {
	client    *ds.Client
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

	if opts.ServiceAccount != nil {
		var conf *oauthjwt.Config
		conf, err = google.JWTConfigFromJSON(
			opts.ServiceAccount,
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

	db.client, err = ds.NewClient(ctx, opts.ProjectID, option.WithTokenSource(ts))
	if err != nil {
		panic(err)
	}

	db.initRole()
	db.initUser()

	return &db
}
