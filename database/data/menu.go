package data

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetInitialMenu() []interface{} {
	return []interface{}{
		bson.M{"name": "Shawai", "price": 550.0, "url": "https://ik.imagekit.io/akhildev/flamespot/shawai.png", "category": "Mains", "created_at": time.Now()},
		bson.M{"name": "Alfaham", "price": 480.0, "url": "https://ik.imagekit.io/akhildev/flamespot/alfaham.png", "category": "Mains", "created_at": time.Now()},
		bson.M{"name": "Kuboos", "price": 15.0, "url": "https://ik.imagekit.io/akhildev/flamespot/kuboos.png", "category": "Breads", "created_at": time.Now()},
		bson.M{"name": "Parotta", "price": 20.0, "url": "https://ik.imagekit.io/akhildev/flamespot/parotta.png", "category": "Breads", "created_at": time.Now()},
		bson.M{"name": "Lime Juice", "price": 25.0, "url": "https://ik.imagekit.io/akhildev/flamespot/lime.png", "category": "Drinks", "created_at": time.Now()},
	}
}
