package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func execsRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", handlers.SignupHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)

	//Google endpoints
	mux.HandleFunc("GET /auth/google/login", handlers.OAuthGoogleLogin)
	mux.HandleFunc("GET /auth/google/callback", handlers.OAuthGoogleCallback)

	// GitHub OAuth2 endpoints
	mux.HandleFunc("GET /auth/github/login", handlers.OAuthGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", handlers.OAuthGitHubCallback)
	return mux
}
