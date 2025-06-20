package handlers

import (
	"traverse/internal/db/redis/cache"
	"traverse/internal/services"

	"github.com/go-playground/validator/v10"
)

type HandlerDeps struct {
	Service *services.Service
	Validator *validator.Validate
	Cache cache.Cache
}

type Handlers struct {
	HealthHandler
	AuthHandler
	UsersHandler
	ContractHandler
	ReviewHandler
}

// TODO: i dont like this constructor for handlers
func New(deps *HandlerDeps) *Handlers {
	return &Handlers{
		NewHealthHandler(),
		NewAuthHandler(deps.Service.Users, deps.Validator, deps.Cache),
		NewUserHandler(deps.Service.Users),
		NewContract(deps.Service.Contract, deps.Validator),
		NewReviewHandler(deps.Service.Review),
	}
}
