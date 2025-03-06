package services

import (
	"traverse/internal/auth"
	"traverse/internal/storage"
)

type Service struct {
	Users UserService
}

func NewServices(storage *storage.Storage, auth auth.Authenticator) *Service {
	return &Service{
		Users: NewUserService(storage, auth),
	}
}
