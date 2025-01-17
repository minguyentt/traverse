package users

import (
	"context"
	"log/slog"
	"net/http"
)

// NOTE: handler => service => repo => database

// collection of user handler interface
type UsersHandler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type usersHandler struct {
	ctx    context.Context
	logger *slog.Logger
}

func NewUsersHandler() *usersHandler {
	logger := slog.Default().With("area", "Users API handler")
	return &usersHandler{
		logger: logger,
	}
}

// pretending this is acting as all the layers...
func (h *usersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("you got all duh users mang"))
}
