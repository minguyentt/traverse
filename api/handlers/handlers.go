package handlers

import (
	"log/slog"
	"github.com/minguyentt/traverse/internal/db/redis/cache"
	"github.com/minguyentt/traverse/internal/services"

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
}

// TODO: i dont like this constructor for handlers
func New(deps *HandlerDeps, l *slog.Logger) *Handlers {
	return &Handlers{
		NewHealthHandler(),
		NewAuthHandler(deps.Service.Users, deps.Validator, deps.Cache, l.With("area", "auth")),
		NewUserHandler(deps.Service.Users),
		NewContract(deps.Service.Contract, deps.Validator, deps.Cache, l.With("area", "contract")),
	}
}
