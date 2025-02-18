package models

/*

NOTE:
1. authenticate the user
2. user then calls getUserByID, parse the urlparam retrieving the requested userID
3. then service will execute passing the parsed userID
4. service communicates with the storage (db)
5. db returns the model obj
6. then transform the obj from db to the json format obj
7. send the response obj back to the client/user

*/

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=6,max=128"`
	Email    string `json:"email"    validate:"required,email,max=255"`
}

type User struct {
	ID            int64  `json:"id"`
	Firstname     string `json:"firstname"`
	Username      string `json:"username"`
	Password      string `json:"-"`
	Email         string `json:"email"`
	IsActive      bool   `json:"is_active"`
	AccountType   AccountType `json:"account_type"`
    CreatedAt     string `json:"created_at"`
}
