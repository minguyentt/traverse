package ctx

import (
	"context"
	"net/http"
	"traverse/models"
)

type userKey struct{}

type feedKey struct{}

type contractKey struct{}

func SetUser(r *http.Request, user any) context.Context {
	return context.WithValue(r.Context(), userKey{}, user)
}

func GetUserFromCTX(r *http.Request) *models.User {
	return r.Context().Value(userKey{}).(*models.User)
}

func SetContract(r *http.Request, user any) context.Context {
	return context.WithValue(r.Context(), userKey{}, user)
}

func GetContractFromCTX(r *http.Request) *models.Contract {
	return r.Context().Value(contractKey{}).(*models.Contract)
}

func SetGlobalFeed(r *http.Request, user any) context.Context {
	return context.WithValue(r.Context(), feedKey{}, user)
}

func GetGlobalFeedFromCTX(r *http.Request) []models.ContractMetaData {
	return r.Context().Value(contractKey{}).([]models.ContractMetaData)
}
