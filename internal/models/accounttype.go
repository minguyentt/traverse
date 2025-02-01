package models

import (
	"traverse/internal/db"
)

type AccountType struct {
	Id          int64  `json:"id"`
	Alias       string `json:"alias"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type AccountTypeStore struct {
	db *db.PGDB
}

// func (s *AccountTypeStore) GetAccountTypeAlias(
// 	ctx context.Context,
// 	alias string,
// ) (*AccountType, error) {
// }
