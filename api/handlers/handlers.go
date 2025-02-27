package handlers

import (
	"github.com/go-playground/validator/v10"
	"traverse/internal/services"
)

// var Valid *validator.Validate
//
// func init() {
//     Valid = validator.New(validator.WithRequiredStructEnabled())
// }

type Handlers struct {
	UsersHandler
	HealthHandler
    ActivateHandler
    AuthHandler
    *validator.Validate
}

// FIX: i dont like this constructor
func NewHandlers(service *services.Service) *Handlers {
    validator := validator.New(validator.WithRequiredStructEnabled())
	return &Handlers{
        &usersHandler{service.Users},
        &healthHandler{},
        &activateUser{service.Activate},
        &authHandler{service.Users, validator},
        validator,
	}
}
