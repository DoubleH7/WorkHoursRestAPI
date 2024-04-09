package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrInvalidAuth = errors.New("incorrect username or password")

func GetAllUsers(client *mongo.Client) (cur *mongo.Cursor, err error) {

	opt := options.Find()
	opt.SetProjection(bson.D{
		{Key: "name", Value: 1},
		{Key: "username", Value: 1},
		{Key: "_id", Value: 1},
	})

	cur, err = client.Database("presenceLog").Collection("users").Find(
		context.TODO(),
		bson.D{{}},
		opt,
	)
	if err != nil {
		return new(mongo.Cursor), err
	}

	return cur, nil
}

func InsertUser(client *mongo.Client, user User) (id string, err error) {
	res01 := client.Database("presenceLog").Collection("users").FindOne(
		context.TODO(),
		bson.D{{Key: "username", Value: user.Username}},
	)
	if res01.Err() != mongo.ErrNoDocuments {
		return "", errors.New("username already exists")
	}

	res02, err := client.Database("presenceLog").Collection("users").InsertOne(
		context.TODO(),
		user,
	)
	if err != nil {
		return "", err
	}
	if res02.InsertedID == new(primitive.ObjectID) {
		return "", errors.New("new user not inserted properly")
	}
	return res02.InsertedID.(primitive.ObjectID).String(), nil
}

func AuthUser(client *mongo.Client, username, password string) (UserWid, error) {
	var user UserWid

	opt := options.FindOne()
	// opt.SetProjection(bson.D{
	// 	{"password", 0},
	// })

	err := client.Database("presenceLog").Collection("users").FindOne(
		context.TODO(),
		bson.D{{Key: "username", Value: username}},
		opt,
	).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return UserWid{}, ErrInvalidAuth
	}
	if err != nil {
		return UserWid{}, err
	}
	if user.Password != password {
		return UserWid{}, ErrInvalidAuth
	}
	user.Password = "**********"

	return user, nil
}
