package model

import "go.mongodb.org/mongo-driver/v2/bson"

type MenuItem struct {
	ID       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Price    float64       `json:"price" bson:"price"`
	URL      string        `json:"url" bson:"url"`
	Category string        `json:"category" bson:"category"`
}
