package cors

import (
	"log"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	log.Println("cors start")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr)
		log.Println(r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{
	// 		"http://127.0.0.1:8001",
	// 		"https://127.0.0.1:8001",
	// 		"http://localhost:8001",
	// 		"https://localhost:8001",
	// 		"http://213.87.152.18:8001",
	// 		"https://213.87.152.18:8001",
	// 		"http://212.233.98.59:8080",
	// 		"https://212.233.98.59:8080"},
	// 	AllowCredentials:   true,
	// 	AllowedMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
	// 	AllowedHeaders:     []string{"*"},
	// 	OptionsPassthrough: false,
	// })
	// return c.Handler(next)
}
