package model

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID       bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"-" bson:"password"`
	Role     string        `json:"role" bson:"role"`
}
