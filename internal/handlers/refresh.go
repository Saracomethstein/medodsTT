package handlers

import (
	"medodsTT/internal/models"
	"medodsTT/internal/services"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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
	request := new(models.RefreshResponse)

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(request.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(services.GetJWTKey()), nil
	})
	if err != nil || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Access token is invalid or expired")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	ip, ok := claims["ip"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims: missing IP")
	}

	refreshTokenHash := h.RefreshService.GetRefreshTokenHash(userID)
	err = bcrypt.CompareHashAndPassword([]byte(refreshTokenHash), []byte(request.RefreshToken))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
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
