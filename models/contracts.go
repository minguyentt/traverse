package models

import "time"

type ContractMetaData struct {
	Contract
	ReviewCounts int `json:"review_counts"`
}

type Contract struct {
	// ID of the contract
	ID int64 `json:"id"`

	// userID tied to the created contract
	UserID  int64    `json:"user_id"`
	Title   string   `json:"name"`
	Address string   `json:"address"`
	City    string   `json:"city"`
	Agency  string   `json:"agency"`
	Reviews []Review `json:"reviews"`
	User    User     `json:"user"`
	// Rating int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}
