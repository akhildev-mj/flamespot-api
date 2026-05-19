package database

import (
	"flamespot-api/config"
	"flamespot-api/database/data"
	"flamespot-api/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RunMigrationsAndSeed(db *mongo.Database) {
	seedMenu(db)
	createIndexes(db)
}

func seedMenu(db *mongo.Database) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := db.Collection(config.CollectionMenu)

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil || count > 0 {
		return
	}

	menuItems := data.GetInitialMenu()
	_, _ = collection.InsertMany(ctx, menuItems)
}

func createIndexes(db *mongo.Database) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "timestamp", Value: -1}},
	}

	_, _ = db.Collection(config.CollectionOrders).Indexes().CreateOne(ctx, indexModel)
}
