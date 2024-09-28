package cors

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func CorsMiddleware(next http.Handler) http.Handler {
	log.Println("cors start")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1",
			"https://127.0.0.1",
			"http://localhost:8001",
			"https://localhost:8001",
			"http://213.87.134.168:8001",
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
