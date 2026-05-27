package model

import "go.mongodb.org/mongo-driver/v2/bson"

type MenuItem struct {
	ID         bson.ObjectID `json:"id" bson:"_id"`
	CategoryID bson.ObjectID `json:"categoryId" bson:"categoryId"`
	Name       string        `json:"name" bson:"name"`
	Price      float64       `json:"price" bson:"price"`
	URL        string        `json:"url" bson:"url"`
	IsActive   bool          `json:"isActive" bson:"isActive"`
	CreatedAt  int64         `json:"-" bson:"created_at"`
	UpdatedAt  int64         `json:"-" bson:"updated_at"`
}
