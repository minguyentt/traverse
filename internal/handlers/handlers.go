package handlers

import (
	"context"
	"traverse/internal/handlers/users"
)

// holds all the HTTP handlers
type APIHandlers struct {
	Ctx   context.Context
	Users users.UsersHandler
}

// TODO: will need a struct to hold all the dependency injections
// db
// cache

func NewHandlers() *APIHandlers {
	return &APIHandlers{
		Users: users.NewUsersHandler(),
	}
}
