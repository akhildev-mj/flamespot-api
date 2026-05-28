package database

import (
	"context"

	"flamespot-api/src/config"

	"go.mongodb.org/mongo-driver/v2/bson"
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

	collections := []string{
		config.CollectionUsers,
		config.CollectionCategories,
		config.CollectionMenu,
		config.CollectionOrders,
	}

	existingNames, err := db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return nil, nil, err
	}

	existingMap := make(map[string]bool)
	for _, name := range existingNames {
		existingMap[name] = true
	}

	for _, collName := range collections {
		if !existingMap[collName] {
			if err := db.CreateCollection(context.TODO(), collName); err != nil {
				return nil, nil, err
			}
		}
	}

	return client, db, nil
}
