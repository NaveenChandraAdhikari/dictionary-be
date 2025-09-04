package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userId int, username, role string) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiresin := os.Getenv("JWT_EXPIRES_IN")

	//create instance of claims that jwt needs
	claims := jwt.MapClaims{
		"uid":  userId,
		"user": username,
		"role": role,
	}

	if jwtExpiresin != "" {

		duration, err := time.ParseDuration(jwtExpiresin)
		if err != nil {
			return "", ErrorHandler(err, "Internal error")
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
		//claims also include the expiry timing for this token
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}
	//create the new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//signed token signed with our secret key
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", ErrorHandler(err, "Internal error ")
	}
	return signedToken, nil

}
