package configs

import (
	"github.com/go-chi/cors"
)

func WithCorsOpts() cors.Options {
	opts := cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
		}, // change to production but local testing for now
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
        AllowCredentials: false,
        MaxAge: 300,
	}

    return opts
}
