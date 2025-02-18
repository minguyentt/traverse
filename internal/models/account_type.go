package models

type AccountType struct {
	ID          int64  `json:"id"`
	AType       string `json:"_type"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}
