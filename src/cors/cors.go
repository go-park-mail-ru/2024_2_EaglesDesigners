package cors

import (
	"net/http"

	"github.com/rs/cors"
)

func CorsMiddleware(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1",
			"https://127.0.0.1",
			"http://localhost",
			"https://localhost",
			"http://213.87.134.168",
			"https://213.87.134.168",
			"http://212.233.98.59",
			"https://212.233.98.59"},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})

	return c.Handler(next)
}
