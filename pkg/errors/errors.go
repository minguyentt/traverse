package errors

import (
	"log/slog"
	"net/http"
	"github.com/minguyentt/traverse/pkg/response"
)

func InternalServerErr(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(
		"Internal server error",
		"METHOD",
		r.Method,
		"PATH",
		r.URL.Path,
		"ERROR",
		err.Error(),
	)

	response.Error(w, http.StatusInternalServerError, "encountered internal server error")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("bad HTTP request", "METHOD", r.Method, "PATH", r.URL.Path, "ERROR", err.Error())

	response.Error(w, http.StatusBadRequest, err.Error())
}

func NotFoundRequest(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(
		"Not found request encounter",
		"method",
		r.Method,
		"path",
		r.URL.Path,
		"error",
		err.Error(),
	)

	response.Error(w, http.StatusNotFound, err.Error())
}

func UnauthorizedBasicErr(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn(
		"Unauthorized basic auth error",
		"METHOD",
		r.Method,
		"PATH",
		r.URL.Path,
		"ERROR",
		err.Error(),
	)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	response.Error(w, http.StatusUnauthorized, "unauthorized")
}

func UnauthorizedErr(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn(
		"Unauthorized auth encounter",
		"METHOD",
		r.Method,
		"PATH",
		r.URL.Path,
		"ERROR",
		err.Error(),
	)

	response.Error(w, http.StatusUnauthorized, "unauthorized")
}
