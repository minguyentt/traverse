package response

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

func Write(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func Read(w http.ResponseWriter, r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func Error(w http.ResponseWriter, status int, msg string) error {
	return Write(w, status, &jsonErr{Error: msg})
}

func JSON(w http.ResponseWriter, status int, data any) error {
	return Write(w, status, &jsonResponse{Data: data})
}
