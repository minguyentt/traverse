package models

/*
Sending => server

- creating user
- consider sending back to client ex. GetUserByID
- what do you want the client to see?
*/

// start with modeling requesting user id object
type User struct {
	Id            int64       `json:"id"`
	Firstname     string      `json:"firstname"`
	Username      string      `json:"username"`
	Email         string      `json:"email"`
	IsActive      bool        `json:"is_active"`
	AccountTypeID string      `json:"account_type_id"`
	AccountType   AccountType `json:"account_type"`

	CreatedAt     string      `json:"created_at"`
}
