package handlers

import (
	"medodsTT/internal/services"

	"github.com/labstack/echo/v4"
)

type GenerateHandler struct {
	GenerateService *services.GenerateService
}

func NewGenerateHandler(generateService *services.GenerateService) *GenerateHandler {
	return &GenerateHandler{GenerateService: generateService}
}

func (h *GenerateHandler) GenerateToken(c echo.Context) error {
	return nil
}
