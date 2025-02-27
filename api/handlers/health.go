package handlers

import "net/http"

type HealthHandler interface {
	HealthChecker(w http.ResponseWriter, r *http.Request)
}

type healthHandler struct{}

func (h *healthHandler) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
