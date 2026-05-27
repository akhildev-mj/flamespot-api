package model

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	Username  string        `json:"username" bson:"username"`
	NickName  string        `json:"nickName" bson:"nickName"`
	Password  string        `json:"-" bson:"password"`
	CreatedAt int64         `json:"-" bson:"created_at"`
	UpdatedAt int64         `json:"-" bson:"updated_at"`
}
