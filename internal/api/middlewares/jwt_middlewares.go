package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"restapi/internal/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

// some trick
type ContextKey string

const (
	RoleKey      ContextKey = "role"
	ExpiresAtKey ContextKey = "expiresAt"
	UsernameKey  ContextKey = "username"
	UserIDKey    ContextKey = "userId"
)

func JWTMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//fmt.Println(r.Cookie("Bearer"))
		token, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Authorization Header missing", http.StatusUnauthorized)
			return
		}

		//token := cookie.Value // ==== just the JWT
		//fmt.Println("Extracted token:", token)

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {

			//custom errros by us
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			} else if errors.Is(err, jwt.ErrTokenMalformed) {
				http.Error(w, "Token malformed", http.StatusUnauthorized)
				return
			}

			utils.ErrorHandler(err, "") //log ,"" because no one is calling the jwtmiddleware
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
			//log.Fatal(err) stop the server
		}

		if parsedToken.Valid {
			log.Print("Valid JWT")

		} else {
			http.Error(w, "Invalid Login(vague error lol) Token", http.StatusUnauthorized)
			log.Println("Invalid JWT", token)
		}
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {

			http.Error(w, "Invalid Login Token", http.StatusInternalServerError)
			log.Println("Invalid login token", token)
			return
		}

		ctx := context.WithValue(r.Context(), RoleKey, claims["role"])
		ctx = context.WithValue(ctx, ExpiresAtKey, claims["exp"])
		ctx = context.WithValue(ctx, UsernameKey, claims["user"])
		ctx = context.WithValue(ctx, UserIDKey, claims["uid"]) // Use UserIDKey

		fmt.Println("this is the context\n", ctx.Value(ContextKey("userId")))
		//fmt.Println("with context\n", r.WithContext(ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
