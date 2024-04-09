package database

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrDuplicateStart = errors.New("user already has an active session")
var ErrDuplicateStop = errors.New("user has no active sessions")

func sessionActive(client *mongo.Client, id primitive.ObjectID) (bool, error) {
	var session Session

	opt := options.FindOne()
	opt.SetSort(bson.D{{Key: "start", Value: -1}})
	err := client.Database("presenceLog").Collection("sessions").FindOne(
		context.TODO(),
		bson.D{{Key: "owner", Value: id}},
		opt,
	).Decode(&session)

	// fmt.Println(session.Start, "\n", session.Stop)

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if session.Start.After(session.Stop) {
		return true, nil
	}

	return false, nil
}

func CreateSession(client *mongo.Client, id primitive.ObjectID) error {
	active, err := sessionActive(client, id)
	if err != nil {
		return err
	}

	if active {
		return ErrDuplicateStart
	}

	var session = Session{
		Owner: id,
		Start: time.Now(),
	}

	_, err = client.Database("presenceLog").Collection("sessions").InsertOne(
		context.TODO(),
		session,
	)
	if err != nil {
		return err
	}

	return nil
}

func FinishSession(client *mongo.Client, id primitive.ObjectID) (dur time.Duration, err error) {

	active, err := sessionActive(client, id)
	if err != nil {
		return 0, err
	}

	if !active {
		return 0, ErrDuplicateStop
	}

	opt := options.FindOneAndUpdate()

	opt.SetSort(bson.D{{Key: "start", Value: -1}})
	opt.SetProjection(bson.D{{Key: "start", Value: 1}})

	now := time.Now()

	var session Session

	err = client.Database("presenceLog").Collection("sessions").FindOneAndUpdate(
		context.TODO(),
		bson.M{"owner": id},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "stop", Value: now},
			}},
		},
		opt,
	).Decode(&session)

	// fmt.Println(session.Start.String(), "\n", session.Stop.String(), "\n", now.String())

	return now.Sub(session.Start), err

}
