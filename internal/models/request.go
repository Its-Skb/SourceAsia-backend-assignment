package models

type RequestPayload struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}
