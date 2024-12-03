package handlers

import (
	"medodsTT/internal/services"

	"github.com/labstack/echo/v4"
)

type RefreshHandler struct {
	RefreshService *services.RefreshService
}

func NewRefreshHandler(refreshService *services.RefreshService) *RefreshHandler {
	return &RefreshHandler{RefreshService: refreshService}
}

func (h *RefreshHandler) RefreshToken(c echo.Context) error {
	return nil
}
