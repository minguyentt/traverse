package services

import (
	"traverse/internal/auth"
	"traverse/internal/storage"
)

type Service struct {
	Users UserService
	Contract ContractService
}

func New(storage *storage.Storage, auth auth.TokenAuthenticator) *Service {
	return &Service{
		Users: NewUserService(storage, auth),
		Contract: NewContractService(storage),
	}
}
