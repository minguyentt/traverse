package models

type AccountType struct {
	ID          int64  `json:"id"`
	Alias       string `json:"alias"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}
