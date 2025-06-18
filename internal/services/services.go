package services

import (
	"traverse/internal/auth"
	"traverse/internal/storage"
)

type Service struct {
	Users UserService
	Contract ContractService
	Review ReviewService
}

func New(storage *storage.Storage, auth auth.Authenticator) *Service {
	return &Service{
		Users: NewUserService(storage, auth),
		Contract: NewContractService(storage),
		Review: NewReviewService(storage),
	}
}
