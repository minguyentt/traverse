package handlers

import (
	"traverse/internal/services"

	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	Users    UsersHandler
	Health   HealthHandler
	Auth     AuthHandler
}

// FIX: i dont like this constructor
func NewHandlers(service *services.Service, validator *validator.Validate) *Handlers {
	return &Handlers{
        Users: NewUserHandler(service),
        Auth: NewAuthHandler(service, validator),
	}
}
