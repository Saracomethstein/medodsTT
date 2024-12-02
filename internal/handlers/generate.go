package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func GenerateToken(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "generate not implemented yet")
}
