package database

import (
	"context"

	"flamespot-api/src/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Connect(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	bsonOpts := &options.BSONOptions{
		ObjectIDAsHexString: true,
	}

	opts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetServerAPIOptions(serverAPI).
		SetBSONOptions(bsonOpts)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, nil, err
	}

	db := client.Database(cfg.DatabaseName)

	return client, db, nil
}
