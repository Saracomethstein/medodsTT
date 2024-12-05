package handlers

import (
	"medodsTT/internal/models"
	"medodsTT/internal/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RefreshHandler struct {
	RefreshService *services.TokenService
}

func NewRefreshHandler(refreshService *services.TokenService) *RefreshHandler {
	return &RefreshHandler{RefreshService: refreshService}
}

func (h *RefreshHandler) RefreshToken(c echo.Context) error {
	request := new(models.RefreshRequest)

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed")
	}

	userID, ip, err := h.RefreshService.GetClaimsFromJWT(*request)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	err = h.RefreshService.CompareHash(userID, *request)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	newAccessToken, err := h.RefreshService.GenerateAccessToken(userID, ip)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate new access token")
	}
	newRefreshToken := h.RefreshService.GenerateRefreshToken()

	hashedRefreshToken, _ := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err := h.RefreshService.SaveRefreshToken(userID, string(hashedRefreshToken), ip, time.Now().Add(7*24*time.Hour)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save refresh token")
	}

	return c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
