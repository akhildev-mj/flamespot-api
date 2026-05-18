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

	log.Println("Checking database collections for missing seed data...")

	seedItems(ctx, db)
}

func seedItems(ctx context.Context, db *mongo.Database) {
	collection := db.Collection("items")

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		log.Printf("Could not count items: %v", err)
		return
	}

	if count == 0 {
		items := []interface{}{
			bson.M{"name": "Shawai", "price": 550.0, "image": "https://ik.imagekit.io/akhildev/flamespot/shawai.png", "category": "Mains", "created_at": time.Now()},
			bson.M{"name": "Alfaham", "price": 480.0, "image": "https://ik.imagekit.io/akhildev/flamespot/alfaham.png", "category": "Mains", "created_at": time.Now()},
			bson.M{"name": "Kuboos", "price": 15.0, "image": "https://ik.imagekit.io/akhildev/flamespot/kuboos.png", "category": "Breads", "created_at": time.Now()},
			bson.M{"name": "Parotta", "price": 20.0, "image": "https://ik.imagekit.io/akhildev/flamespot/parotta.png", "category": "Breads", "created_at": time.Now()},
			bson.M{"name": "Lime Juice", "price": 25.0, "image": "https://ik.imagekit.io/akhildev/flamespot/lime.png", "category": "Drinks", "created_at": time.Now()},
		}

		if _, err := collection.InsertMany(ctx, items); err != nil {
			log.Printf("Failed to seed items: %v", err)
		} else {
			log.Println("Seeded initial Items successfully!")
		}
	}
}
