package errors

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Data any `json:"data"`
}

type jsonErr struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func jsonWithErr(w http.ResponseWriter, status int, msg string) error {
	return writeJSON(w, status, &jsonErr{Error: msg})
}

func JSONResponse(w http.ResponseWriter, status int, data any) error {
	return writeJSON(w, status, &jsonResponse{Data: data})
}
