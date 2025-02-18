package handlers

import (
	"github.com/minguyentt/traverse/internal/services"
)

type Handlers struct {
	UsersHandler
	HealthHandler
}

func NewHandlers(service *services.Service) *Handlers {
	return &Handlers{
        &usershandler{service.Users},
        &healthHandler{},
	}
}
