package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/pkg/utils"
	"restapi/internal/repository/sqlconnect"
	"time"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	var newExecs models.Exec

	var rawExecs map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ERROR READING REQUEST BODY", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(body, &rawExecs)
	if err != nil {
		log.Println("Error unmarshalling to map")
		http.Error(w, "INVALID REQUEST BODY", http.StatusBadRequest)
		fmt.Println("raw exec", rawExecs)
		return
	}

	fields := GetFieldNames(models.Exec{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for key := range rawExecs {
		_, ok := allowedFields[key]
		if !ok {
			http.Error(w, "Unacceptable filed found in request.Only use allowed fields.", http.StatusBadRequest)
			return
		}
	}

	err = json.Unmarshal(body, &newExecs)
	if err != nil {
		log.Println("Error unmarshalling to struct")
		http.Error(w, "INVALID REQUEST BODY", http.StatusBadRequest)
		return

	}
	fmt.Println("newexecs", newExecs)

	err = CheckBlankFields(newExecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//DB thing
	addedExec, err := sqlconnect.AddExecDBHandler(newExecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//inline struct
	response := struct {
		Status string      `json:"status"`
		Data   models.Exec `json:"data"`
	}{
		Status: "success",
		Data:   addedExec,
	}
	json.NewEncoder(w).Encode(response)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req models.Exec

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.HashedPassword == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	///TODO user not found error where ???
	user, err := sqlconnect.GetUserByUsername(req.Username)
	fmt.Println("user", user)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusNotFound)
		return
	}

	if user.InactiveStatus {
		http.Error(w, "account is inactive", http.StatusForbidden)
		return
	}
	//TODO (dont know deal with hashed_password not null issue ASAP) Check if user is OAuth user
	if user.IsOAuthUser || user.HashedPassword == "NO_PASSWORD_OAUTH_USER" {
		http.Error(w, "This account uses OAuth login. Please use Google/GitHub login.", http.StatusBadRequest)
		return
	}
	err = utils.VerifyPassword(req.HashedPassword, user.HashedPassword)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusNotFound)
		return
	}
	tokenString, err := utils.SignToken(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "coudl not create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{

		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
		//SameSite: http.SameSiteNoneMode
	},
	)
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(response)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}
