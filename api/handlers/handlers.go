package handlers

import (
	"traverse/internal/services"

	"github.com/go-playground/validator/v10"
)

type Handlers struct {
    Auth     AuthHandler
	Users    UsersHandler
	Health   HealthHandler
}

// FIX: i dont like this constructor
func NewHandlers(service *services.Service, validator *validator.Validate) *Handlers {
	return &Handlers{
        Auth: NewAuthHandler(service.Users, validator),
        Users: NewUserHandler(service.Users),
	}
}
