package handlers

import (
	"log/slog"
	"traverse/internal/db/redis/cache"
	"traverse/internal/services"

	"github.com/go-playground/validator/v10"
)

type ctxKeyUser struct{}

type ctxGlobalFeed struct{}

type HandlerDeps struct {
	Service   *services.Service
	Validator *validator.Validate
	Cache     cache.Redis
	Logger    *slog.Logger
}

type Handlers struct {
	HealthHandler
	AuthHandler
	UsersHandler
	ContractHandler
	ReviewHandler
}

// TODO: i dont like this constructor for handlers
func New(deps *HandlerDeps, l *slog.Logger) *Handlers {
	return &Handlers{
		NewHealthHandler(),
		NewAuthHandler(deps.Service.Users, deps.Validator, deps.Cache),
		NewUserHandler(deps.Service.Users),
		NewContract(deps.Service.Contract, deps.Validator, deps.Cache, l),
		NewReviewHandler(deps.Service.Review),
	}
}
