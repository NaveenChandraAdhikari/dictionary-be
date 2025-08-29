package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(password, encodedHash string) error {
	//split stored hash into salt and real hash
	parts := strings.Split(encodedHash, ",")
	fmt.Println(parts)
	if len(parts) != 2 {
		return ErrorHandler(errors.New("invalid encoded hash format"), "internal server error")
		//http.Error(w, "invalid encoded hash format", http.StatusForbidden)
		//return true
	}

	saltBase64 := parts[0]
	hashPasswordBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return ErrorHandler(err, "internal server error")
		//http.Error(w, "failed to decode the salt", http.StatusForbidden)
		//return true
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(hashPasswordBase64)
	if err != nil {
		return ErrorHandler(err, "internal server error")
		//http.Error(w, "failed to decode the hashed password", http.StatusForbidden)
		//return true
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	//to compare the hash it is not like == there are steps first check the len of the passwords
	if len(hash) != len(hashedPassword) {
		return ErrorHandler(errors.New("hash length mismatch"), "incorrect password`")
	}
	if subtle.ConstantTimeCompare(hash, hashedPassword) == 1 {
		return nil
		//do nothing
	}
	return ErrorHandler(errors.New("incorrect password`"), "incorrect password`")

}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrorHandler(errors.New("password is blank"), "please enter password")
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", ErrorHandler(errors.New("failed to generate salt"), "internal error")
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s,%s", saltBase64, hashBase64)
	//password = encodedHash
	return encodedHash, nil
}
