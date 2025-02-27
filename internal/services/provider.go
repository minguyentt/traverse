package services

import (
	"traverse/internal/auth"
	"traverse/internal/storage"
)

type Service struct {
	Users *UserService
    Activate *ActivateService
}

func NewServices(storage *storage.Storage, auth auth.Authenticator) *Service {
	return &Service{
		Users: NewUserService(storage, auth),
        Activate: NewActivationService(storage),
	}
}
