package services

import (
	"github.com/minguyentt/traverse/internal/storage"
)

type Service struct {
	Users *UserService
}

func NewServices(storage *storage.Storage) *Service {
	return &Service{
		Users: NewUserService(storage),
	}
}
