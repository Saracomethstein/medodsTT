package handlers

import (
	"medodsTT/internal/models"
	"medodsTT/internal/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type GenerateHandler struct {
	GenerateService *services.TokenService
}

func NewGenerateHandler(generateService *services.TokenService) *GenerateHandler {
	return &GenerateHandler{GenerateService: generateService}
}

func (h *GenerateHandler) GenerateToken(c echo.Context) error {
	req := new(models.TokenRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, err := h.GenerateService.GenerateAccessToken(req.UserID, req.IP)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken := h.GenerateService.GenerateRefreshToken()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash refresh token")
	}

	err = h.GenerateService.SaveRefreshToken(req.UserID, string(hashedToken), req.IP, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save refresh token")
	}

	return c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
