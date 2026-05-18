package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RunMigrationsAndSeed(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Checking database collections and indexes...")

	// 1. Seed the menu data
	seedMenu(ctx, db)

	// 2. Create required indexes for pagination
	createIndexes(db)
}

func seedMenu(ctx context.Context, db *mongo.Database) {
	collection := db.Collection("menu")

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		log.Printf("Could not count menu items: %v", err)
		return
	}

	if count == 0 {
		menuItems := []interface{}{
			bson.M{"name": "Shawai", "price": 550.0, "url": "https://ik.imagekit.io/akhildev/flamespot/shawai.png", "category": "Mains", "created_at": time.Now()},
			bson.M{"name": "Alfaham", "price": 480.0, "url": "https://ik.imagekit.io/akhildev/flamespot/alfaham.png", "category": "Mains", "created_at": time.Now()},
			bson.M{"name": "Kuboos", "price": 15.0, "url": "https://ik.imagekit.io/akhildev/flamespot/kuboos.png", "category": "Breads", "created_at": time.Now()},
			bson.M{"name": "Parotta", "price": 20.0, "url": "https://ik.imagekit.io/akhildev/flamespot/parotta.png", "category": "Breads", "created_at": time.Now()},
			bson.M{"name": "Lime Juice", "price": 25.0, "url": "https://ik.imagekit.io/akhildev/flamespot/lime.png", "category": "Drinks", "created_at": time.Now()},
		}

		if _, err := collection.InsertMany(ctx, menuItems); err != nil {
			log.Printf("Failed to seed menu: %v", err)
		} else {
			log.Println("Seeded initial Menu successfully!")
		}
	}
}

func createIndexes(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "timestamp", Value: -1}},
	}

	_, err := db.Collection("orders").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Note: Index creation on 'orders.timestamp' skipped or failed: %v", err)
	} else {
		log.Println("Index on 'orders.timestamp' verified/created.")
	}
}
