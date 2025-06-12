package models

type Review struct {
	ID       int64  `json:"id"`
	ReviewID int64  `json:"review_id"`
	UserID   int64  `json:"user_id"`
	Content  string `json:"content"`
	User     User   `json:"user"`
}
