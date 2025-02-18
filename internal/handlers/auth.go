package handlers

import "net/http"

type AuthHandler interface {
	RegisterUser(http.ResponseWriter, *http.Request)
    CreateUserToken(http.ResponseWriter, *http.Request)
}

type authhandler struct {
}



