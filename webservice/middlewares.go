package webservice

import (
	"fmt"
	"net/http"

	"github.com/DoubleH7/presenceHoursLog/database"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserpassCheck(client *mongo.Client) func(string, string, echo.Context) (bool, error) {
	return func(username, password string, c echo.Context) (bool, error) {

		//fetching user from database based on credentials
		user, err := database.AuthUser(client, username, password)

		if err == database.ErrInvalidAuth {
			return false, c.String(http.StatusOK, "incorrect username or password")
		}
		if err != nil {
			fmt.Println("error when authenticating user: ", err)
			return false, c.String(http.StatusInternalServerError, "sorry, something went wrong on our side")
		}

		// saving user in the context for later use
		c.Set("user", user)

		return true, nil
	}
}
