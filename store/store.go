// Package store provides functionality to connect to a
// MongoDB instance and perform database operations.
package store

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	DatabaseName string
	*mongo.Client
}

// For ease of unit testing.
var (
	newClient = func(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
		return mongo.Connect(ctx, opts)
	}
	ping = func(ctx context.Context, client *mongo.Client) error {
		return client.Ping(ctx, nil)
	}
)

// Connect connects to a running MongoDB instance.
func Connect(ctx context.Context, host, database string, port int) (*MongoDb, error) {
	newClient, err := newClient(ctx, options.Client().ApplyURI(
		uri(host, port),
	))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create MongoDB client")
	}
	err = ping(ctx, newClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping MongoDB server")
	}
	return &MongoDb{
		DatabaseName: database,
		Client:       newClient,
	}, nil
}

// uri generates uri string for connecting to MongoDB.
func uri(host string, port int) string {
	const format = "mongodb://%s:%d"
	return fmt.Sprintf(format, host, port)
}
