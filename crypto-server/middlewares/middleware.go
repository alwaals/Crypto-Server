package middlewares

import (
	"log"
	"net/http"
)

func ValidateHeader(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println("Reached Middleware with endpoint details:", req.URL.Path+" , "+req.Method)
		if req.Method != "GET" {
			http.Error(w, "Supports only GET method", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	}
}
