package services

import (
	"github.com/minguyentt/traverse/internal/auth"
	"github.com/minguyentt/traverse/internal/storage"
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
