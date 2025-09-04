package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrigins = []string{
	"http://localhost:3000",
	"http://localhost:44609",
	"http://127.0.0.1:3000",
	"http://127.0.0.1:44609",
	"http://127.0.0.1:36225",
	"http://127.0.0.1:36041",
	"http://127.0.0.1:44867",
}

func Cors(next http.Handler) http.Handler {

	fmt.Println("Cors Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		//
		//if isOriginAllowed(origin) {
		//	w.Header().Set("Access-Control-Allow-Origin", origin)
		//
		//} else {
		//	http.Error(w, "Not allowed by Cors", http.StatusForbidden)
		//}
		if isOriginAllowed(origin) || origin == "" {

			if origin == "" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			http.Error(w, "Not allowed by Cors", http.StatusForbidden)
			return
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
func isOriginAllowed(origin string) bool {

	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
