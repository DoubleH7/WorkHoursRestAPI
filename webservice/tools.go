package webservice

import (
	"net/http"

	"github.com/labstack/echo"
)

func MyInfo(c echo.Context) error {
	user := c.Get("user")
	return c.JSON(http.StatusOK, user)
}

func ServerAlive(c echo.Context) error {
	return c.String(http.StatusOK, "Server is alive and listening")
}
