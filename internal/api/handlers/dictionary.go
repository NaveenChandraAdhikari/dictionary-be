package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"restapi/internal/api/middlewares"
	"restapi/internal/models"
	"restapi/internal/pkg/utils"
	"restapi/internal/repository/sqlconnect"
	"strconv"
)

func AddWordHandler(w http.ResponseWriter, r *http.Request) {

	var newWord models.Word

	var rawWord map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &rawWord)
	if err != nil {
		log.Println("Error unmarshalLing to map")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Word{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for key := range rawWord {
		_, ok := allowedFields[key]
		if !ok {
			http.Error(w, "unacceptable field found in request.Only use allowed fields", http.StatusBadRequest)
			return
		}
	}

	err = json.Unmarshal(body, &newWord)
	if err != nil {
		log.Println("Error unmarshalLing to struct")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return

	}
	fmt.Println("new word", newWord)

	err = CheckBlankFields(newWord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//jwt context so we grab here as per need
	//createdBy := r.Context().Value("user_id")

	//createdBy := r.Context().Value("user_id").(int)
	//fmt.Println("CONTEXT   \n", r.Context())
	//fmt.Println("please print\n", r.Context().Value("userId"))
	//userID, ok := r.Context().Value("userId").(int)
	//if !ok {
	//	http.Error(w, "userid not found in context", http.StatusUnauthorized)
	//	return
	//}
	//createdBy := userID
	// Retrieve userId from context using ContextKey
	userIDValue, ok := r.Context().Value(middlewares.UserIDKey).(float64)
	if !ok {
		http.Error(w, "userId not found in context or invalid type", http.StatusUnauthorized)
		return
	}
	createdBy := int(userIDValue) // Convert float64 to int
	fmt.Println("please print\n", createdBy)

	addedWord, err := sqlconnect.AddWordDBHandler(newWord, createdBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string      `json:"status"`
		Data   models.Word `json:"data"`
	}{
		Status: "success",
		Data:   addedWord,
	}

	json.NewEncoder(w).Encode(response)

}

func DeleteWordHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDValue, ok := r.Context().Value(middlewares.UserIDKey).(float64)
	if !ok {
		http.Error(w, "userId not found in context or invalid type", http.StatusUnauthorized)
		return
	}
	userID := int(userIDValue)

	deletedIds, err := sqlconnect.DeleteWordDBHandler(ids, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIds []int  `json:"deleted_ids"`
	}{
		Status:     "Words successfully deleted",
		DeletedIds: deletedIds,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteOneWordHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userIDValue, ok := r.Context().Value(middlewares.UserIDKey).(float64)
	if !ok {
		http.Error(w, "userId not found in context or invalid type", http.StatusUnauthorized)
		return
	}
	userID := int(userIDValue)

	err = sqlconnect.DeleteOneWordDBHandler(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status    string `json:"status"`
		DeletedId int    `json:"deleted_id"`
	}{
		Status:    "Word successfully deleted",
		DeletedId: id,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateWordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	userIDValue, ok := r.Context().Value(middlewares.UserIDKey).(float64)
	if !ok {
		http.Error(w, "userId not found in context or invalid type", http.StatusUnauthorized)
		return
	}
	userID := int(userIDValue)

	var updatedWord models.Word
	err = json.NewDecoder(r.Body).Decode(&updatedWord)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	updatedWordFromDB, err := sqlconnect.UpdatedWordDBHandler(id, updatedWord, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := struct {
		Status string      `json:"status"`
		Data   models.Word `json:"data"`
	}{
		Status: "Word successfully updated",
		Data:   updatedWordFromDB,
	}

	json.NewEncoder(w).Encode(response)
}
func GetWordsHandler(w http.ResponseWriter, r *http.Request) {

	userIDValue, ok := r.Context().Value(middlewares.UserIDKey).(float64)
	if !ok {
		http.Error(w, "userId not found in context or invalid type", http.StatusUnauthorized)
		return
	}
	userID := int(userIDValue)

	page, limit := utils.GetPaginationParams(r)

	words, totalWords, err := sqlconnect.GetWordsDBHandler(r, userID, limit, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := struct {
		Status   string        `json:"status"`
		Count    int           `json:"count"`
		Page     int           `json:"page"`
		PageSize int           `json:"page_size"`
		Data     []models.Word `json:"data"`
	}{
		Status:   "Words successfully retrieved",
		Count:    totalWords,
		Page:     page,
		PageSize: limit,
		Data:     words,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func GetOneWordHandler(w http.ResponseWriter, r *http.Request) {
	word := r.PathValue("word")

	dbWord, err := sqlconnect.GetOneWordDBHandler(word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbWord)
}
