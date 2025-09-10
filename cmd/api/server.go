package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/middlewares/router"
	"restapi/internal/pkg/oauth"
	"restapi/internal/pkg/utils"
	"restapi/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	oauth.InitOAuthConfigs()

	//connect to postgres
	_, err = sqlconnect.ConnectDb()
	if err != nil {
		utils.ErrorHandler(err, "")
		return
	}

	//get port from env
	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":8080"
	}

	//cert := "cmd/api/cert.pem"
	//key := "cmd/api/key.pem"
	//
	//tlsConfig := &tls.Config{
	//	MinVersion: tls.VersionTLS12,
	//}
	router := router.MainRouter()

	//jwtMiddleware := mw.MiddlewaresExcludePaths(mw.JWTMiddleware, "/signup", "/login")
	//
	//secureMux := mw.Cors(jwtMiddleware(mw.SecurityHeaders(router)))
	//
	////secureMux := mw.SecurityHeaders(router)
	//server := &http.Server{
	//	Addr:    port,
	//	Handler: secureMux,
	//	//TLSConfig: tlsConfig,
	//}
	// Exclude JWT for public routes like signup/login
	jwtMiddleware := mw.MiddlewaresExcludePaths(mw.JWTMiddleware, "/signup", "/login", "/auth/google/callback", "/auth/github/callback", "/auth/google/login", "/auth/github/login")

	secureMux := mw.Cors(jwtMiddleware(mw.SecurityHeaders(router)))

	server := &http.Server{
		Addr:    port,
		Handler: secureMux,
	}
	fmt.Println("Starting server at port:", port)

	//err = server.ListenAndServeTLS(cert, key)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
