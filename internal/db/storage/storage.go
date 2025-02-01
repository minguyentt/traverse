package storage

import (
	"context"
	"traverse/internal/models"
)

type Storage struct {
	Users interface {
		CreateUser(context.Context, *models.User) error
		GetUserByID(context.Context, int64) error
		DeleteUser(context.Context, int64) error
	}
}
