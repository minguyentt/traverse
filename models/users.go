package models

import (
	"time"
	"traverse/internal/auth"
)

type User struct {
	ID        int64         `json:"id"`
	Firstname string        `json:"firstname"`
	Username  string        `json:"username"`
	Password  auth.Password `json:"-"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"created_at"`
}

type RegistrationPayload struct {
	Firstname string `json:"firstname" validate:"required,min=1,max=50"`
	Username  string `json:"username"  validate:"required,min=5,max=50"`
	Password  string `json:"password"  validate:"required,min=6,max=128"`
	Email     string `json:"email"     validate:"required,email,max=255"`
}

type UserToken struct {
	Token string `json:"token"`
}

// creating user token
type UserTokenPayload struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=6,max=128"`
	Email    string `json:"email"    validate:"required,email,max=255"`
}

type UserLoginPayload struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=6,max=128"`
}

// type AccountType struct {
// 	ID          int64  `json:"id"`
// 	AType       string `json:"_type"`
// 	Level       int    `json:"level"`
// 	Description string `json:"description"`
// }
