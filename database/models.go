package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Owner primitive.ObjectID `json:"owner" bson:"owner"`
	Start time.Time          `json:"start" bson:"start"`
	Stop  time.Time          `json:"stop" bson:"stop"`
}

type User struct {
	FullName string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Access   uint8  `json:"access" bson:"access"`
	Age      uint8  `json:"age" bson:"age"`
}

type UserWid struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	FullName string             `json:"name" bson:"name"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Access   uint8              `json:"access" bson:"access"`
	Age      uint8              `json:"age" bson:"age"`
}

var AccessLevels = map[string]int{
	"supreme": 3,
	"admin":   2,
	"basic":   1,
	"none":    0,
}
