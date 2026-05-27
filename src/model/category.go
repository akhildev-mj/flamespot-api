package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Category struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	URL       string        `json:"url" bson:"url"`
	IsActive  bool          `json:"isActive" bson:"isActive"`
	CreatedAt int64         `json:"-" bson:"created_at"`
	UpdatedAt int64         `json:"-" bson:"updated_at"`
}
