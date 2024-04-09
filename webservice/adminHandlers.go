package webservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DoubleH7/presenceHoursLog/database"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {

		cur, err := database.GetAllUsers(client)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Sorry something went wrong on our side")
		}

		var results []struct {
			FullName string             `json:"name" bson:"name"`
			Username string             `json:"username" bson:"username"`
			ID       primitive.ObjectID `json:"id" bson:"_id"`
		}

		err = cur.All(context.TODO(), &results)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Sorry something went wrong on our side")
		}

		return c.JSON(http.StatusOK, results)
	}
}

func CreateUser(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		datajs, err := io.ReadAll(c.Request().Body)
		defer c.Request().Body.Close()

		if err != nil {
			return c.String(http.StatusBadRequest, "failed reading request data")
		}

		var user database.User
		err = json.Unmarshal(datajs, &user)

		if err != nil {
			return c.String(http.StatusBadRequest, "failed unmarshalling request data")
		}

		var zUser database.User
		if user.FullName == zUser.FullName ||
			user.Username == zUser.Username ||
			user.Password == zUser.Password {
			return c.String(http.StatusBadRequest, "provide data in the proper format\nFull name, username and password are neccessary")
		}

		id, err := database.InsertUser(client, user)

		if err != nil {
			if err.Error() == "username already exists" {
				return c.String(http.StatusNotAcceptable, fmt.Sprintf("username %s is has already been taken", user.Username))
			}
			return c.String(http.StatusInternalServerError, "Something went wrong on our side")
		}

		return c.String(http.StatusOK, fmt.Sprintf("New user successfully created with the ID : %s", id))
	}
}
