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

	fmt.Println("-----------------------------JWT middleware-----------------------")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("-----------------------------Inside JWT middleware-----------------------")
		//fmt.Println(r.Cookie("Bearer"))
		token, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Authorization Header missing", http.StatusUnauthorized)
			return
		}

		//token := cookie.Value // <-- just the JWT
		//fmt.Println("Extracted token:", token)

		jwtSecret := os.Getenv("JWT_SECRET")

		//MOST IMPORTANT PART ;;;;--parsing and validating the token
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

		//fmt.Println("Parsed token", parsedToken)
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		//if ok {
		//	fmt.Println(claims["uid"], claims["exp"], claims["role"])
		//
		//} else
		if !ok {

			http.Error(w, "Invalid Login Token", http.StatusInternalServerError)
			log.Println("Invalid login token", token)
			return
		}
		//
		//ctx := context.WithValue(r.Context(), "role", claims["role"])
		//ctx = context.WithValue(ctx, "expiresAt", claims["exp"])
		//ctx = context.WithValue(ctx, "username", claims["user"])
		//ctx = context.WithValue(ctx, "userId", claims["uid"])
		//WE USED HERE THE VALUE AS NOT STRING TYPES WWE USE AS CONTEXT KEY TYPES ,,JUST A TRICK not so ,,it is just a name ,i think go compier doesnt want us to use explicitly string types ,,so we use implcitly
		//ctx := context.WithValue(r.Context(), ContextKey("role"), claims["role"])
		//ctx = context.WithValue(ctx, ContextKey("expiresAt"), claims["exp"])
		//ctx = context.WithValue(ctx, ContextKey("username"), claims["user"])
		//ctx = context.WithValue(ctx, ContextKey("userId"), claims["uid"])

		ctx := context.WithValue(r.Context(), RoleKey, claims["role"])
		ctx = context.WithValue(ctx, ExpiresAtKey, claims["exp"])
		ctx = context.WithValue(ctx, UsernameKey, claims["user"])
		ctx = context.WithValue(ctx, UserIDKey, claims["uid"]) // Use UserIDKey

		fmt.Println("this is the context\n", ctx.Value(ContextKey("userId")))
		//fmt.Println("with context\n", r.WithContext(ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
		fmt.Println("-----------------------------Sent response from JWT middleware-----------------------")
	})

}

//we are going to pass the value from function to function internally inside our apis that to be done using context
//use context to carry the claims information across different middlewares or functions
