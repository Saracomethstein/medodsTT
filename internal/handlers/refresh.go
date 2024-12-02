package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func RefreshToken(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "refresh not implemented yet")
}
