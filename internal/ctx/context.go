package ctx

import (
	"context"
	"net/http"
	"traverse/models"
)

func SetUser(r *http.Request, key string, val any) context.Context {
	return context.WithValue(r.Context(), key, val)
}

func GetUserFromCTX(r *http.Request, key string) *models.User {
	return r.Context().Value(key).(*models.User)
}

func SetContract(r *http.Request, key string, val any) context.Context {
	return context.WithValue(r.Context(), key, val)
}

func GetContractFromCTX(r *http.Request, key string) *models.Contract {
	return r.Context().Value(key).(*models.Contract)
}

func SetGlobalFeed(r *http.Request, key string, val any) context.Context {
	return context.WithValue(r.Context(), key, val)
}

func GetGlobalFeedFromCTX(r *http.Request, key string) []models.ContractMetaData {
	return r.Context().Value(key).([]models.ContractMetaData)
}
