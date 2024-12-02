package models

type Request struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
}
