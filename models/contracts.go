package models

import "time"

type Contract struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Location     Location `json:"location"`
	Agency       string   `json:"agency"`
	Reviews      []Review `json:"reviews"`
	ReviewCounts int      `json:"review_counts"`
    // Rating int
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

type Location struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	RegionCode int    `json:"region_code"`
}
