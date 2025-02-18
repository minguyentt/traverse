package errors

import (
	"log/slog"
	"net/http"
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

	jsonWithErr(w, http.StatusInternalServerError, "encountered internal server error")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("bad HTTP request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	jsonWithErr(w, http.StatusBadRequest, err.Error())
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

	jsonWithErr(w, http.StatusNotFound, err.Error())
}
