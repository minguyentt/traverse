package valids

import (
	"github.com/go-playground/validator/v10"
)

// checking valid struct payloads for requests/responses
type valid struct {
	*validator.Validate
}

func New() *valid {
	return &valid{validator.New(validator.WithRequiredStructEnabled())}
}
