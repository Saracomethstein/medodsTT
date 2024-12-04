package handlers

import (
	"log"
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
		return []byte("J8sK^7z!fA0p@o3wY%M#E1Qx%Rk4U&Nv2KZ"), nil
	})
	if err != nil || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Access token is invalid or expired")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	// i check refresh token with bcrypt refesh token (bad idea)
	refreshTokenHash := h.RefreshService.GetRefreshTokenHash(userID)
	err = bcrypt.CompareHashAndPassword([]byte(refreshTokenHash), []byte(request.RefreshToken))
	if err != nil {
		log.Println("Token mismatch:", err, "Token in request: ", request.RefreshToken, "Token in data base: ", refreshTokenHash)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	// for generated new access token i need userID and ip (new function)
	newAccessToken, err := h.RefreshService.GenerateAccessToken(userID, "")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate new access token")
	}

	newRefreshToken := h.RefreshService.GenerateRefreshToken()

	hashedRefreshToken, _ := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err := h.RefreshService.SaveRefreshToken(userID, string(hashedRefreshToken), c.RealIP(), time.Now().Add(7*24*time.Hour)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save refresh token")
	}

	return c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
