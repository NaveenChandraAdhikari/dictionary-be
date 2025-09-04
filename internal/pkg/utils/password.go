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
	// we have the actualy hashed of the passwrod that the user used while user signed up ,,now create another hash of the password the user log in IN,,
	// now compare this hash to the existing hashedPassword we have
	//hash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	//to compare the hash it is not like == there are steps first check the len of the passwords
	if len(hash) != len(hashedPassword) {
		return ErrorHandler(errors.New("hash length mismatch"), "incorrect password`")
		//http.Error(w, "incorrect password", http.StatusForbidden)
		//return true
	}
	//if the length is same then we compare hashes
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
	//hashong and convert string to byte slice
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	//again encode and also we  dont directly save it
	//first encode the salt
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	//geenrate base64 value for hash
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s,%s", saltBase64, hashBase64)
	//password = encodedHash
	return encodedHash, nil
}
