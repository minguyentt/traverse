package errors

import (
	"log/slog"
	"net/http"

	"traverse/api/json"
)

func InternalServerErr(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(
		"internal error",
		"method",
		r.Method,
		"path",
		r.URL.Path,
		"error",
		err.Error(),
	)

	json.ErrResponse(w, http.StatusInternalServerError, "encountered internal server error")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("bad HTTP request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	json.ErrResponse(w, http.StatusBadRequest, err.Error())
}

func NotFoundRequest(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(
		"not found request",
		"method",
		r.Method,
		"path",
		r.URL.Path,
		"error",
		err.Error(),
	)

	json.ErrResponse(w, http.StatusNotFound, err.Error())
}

func UnauthorizedBasicErr(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn(
		"unauthorized basic auth error",
		"method",
		r.Method,
		"path",
		r.URL.Path,
		"error",
		err.Error(),
	)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	json.ErrResponse(w, http.StatusUnauthorized, "unauthorized")
}

func UnauthorizedErr(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn(
		"unauthorized auth encounter",
		"method",
		r.Method,
		"path",
		r.URL.Path,
		"error",
		err.Error(),
	)

	json.ErrResponse(w, http.StatusUnauthorized, "unauthorized")
}
