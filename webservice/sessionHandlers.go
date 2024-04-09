package webservice

import (
	"fmt"
	"net/http"

	"github.com/DoubleH7/presenceHoursLog/database"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartSession(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {

		user := c.Get("user").(database.UserWid)
		if user.Access < uint8(database.AccessLevels["basic"]) {
			return c.String(http.StatusUnauthorized, "you don't have access to this feature.")
		}

		err := database.CreateSession(client, user.ID)
		if err == database.ErrDuplicateStart {
			return c.String(http.StatusBadRequest, "You already have an active session.")
		}
		if err != nil {
			return c.String(http.StatusInternalServerError, "Sorry! Something went wrong on our side.")
		}

		return c.String(http.StatusOK, fmt.Sprintf("A new Session was started for %s.", user.FullName))

	}
}

func StopSession(client *mongo.Client) func(c echo.Context) error {
	return func(c echo.Context) error {

		user := c.Get("user").(database.UserWid)
		if user.Access < uint8(database.AccessLevels["basic"]) {
			return c.String(http.StatusUnauthorized, "you don't have access to this feature.")
		}

		dur, err := database.FinishSession(client, user.ID)
		if err == database.ErrDuplicateStop {
			return c.String(http.StatusBadRequest, "You have no active sessions to Finish")
		}
		if err != nil {
			return c.String(http.StatusInternalServerError, "Sorry! Something went wrong on our side.")
		}

		return c.String(http.StatusOK, fmt.Sprintf("%s's session ended after %dh %dm %ds.",
			user.FullName,
			int(dur.Hours()),
			int(dur.Minutes())%60,
			int(dur.Seconds())%60,
		))

	}
}
