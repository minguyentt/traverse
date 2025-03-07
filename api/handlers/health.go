package handlers

import "net/http"

type HealthHandler interface {
	HealthChecker(w http.ResponseWriter, r *http.Request)
}

type healthhandler struct {}
// TODO: properly add more to ping db too? idk
func NewHealthHandler() *healthhandler {
    return &healthhandler{}
}
func (h *healthhandler) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
