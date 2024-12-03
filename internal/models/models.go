package models

type TokenRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	IP     string `json:"ip" validate:"required,ip"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
