package handlers

import (
	"traverse/internal/services"

	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	HealthHandler
	AuthHandler
	UsersHandler
	ContractHandler
	ReviewHandler
}

// TODO: i dont like this constructor for handlers
func New(service *services.Service, validator *validator.Validate) *Handlers {
	return &Handlers{
		NewHealthHandler(),
		NewAuthHandler(service.Users, validator),
		NewUserHandler(service.Users),
		NewContract(service.Contract, validator),
		NewReviewHandler(service.Review),
	}
}
