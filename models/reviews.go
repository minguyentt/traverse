package models

type Review struct {
	ID         int64  `json:"id"`
	ContractID int64  `json:"contract_id"`
	UserID     int64  `json:"user_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	User       User   `json:"user"`
}
