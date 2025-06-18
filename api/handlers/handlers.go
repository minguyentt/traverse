package handlers

import (
	"net/http"
	"traverse/internal/services"
	"traverse/models"

	"github.com/go-playground/validator/v10"
)

type UserKey string

type Handlers struct {
	HealthHandler
	AuthHandler
	UsersHandler
	ContractHandler
}

// TODO: i dont like this constructor for handlers
func New(service *services.Service, validator *validator.Validate) *Handlers {
	return &Handlers{
		NewHealthHandler(),
		NewAuthHandler(service.Users, validator),
		NewUserHandler(service.Users),
		NewContract(service.Contract, validator),
	}
}

func GetUserCtx(r *http.Request) *models.User {
	usr := r.Context().Value("user").(*models.User)
	return usr
}
