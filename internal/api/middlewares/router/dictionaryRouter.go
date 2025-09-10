package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func dictionaryRouter() *http.ServeMux {
	mux := http.NewServeMux()

	//add a word
	mux.HandleFunc("POST /words/", handlers.AddWordHandler)
	//update existing word
	mux.HandleFunc("PUT /words/{id}", handlers.UpdateWordHandler)
	////fetch a  words meaning
	mux.HandleFunc("GET /words/", handlers.GetWordsHandler)
	//mux.HandleFunc("GET /words", handlers.GetWordsHandler)
	mux.HandleFunc("GET /word/{word}", handlers.GetOneWordHandler)
	////list with pagination, search , sort
	//mux.HandleFunc("GET /words", handlers.ListWordsHandler)
	////delete
	mux.HandleFunc("DELETE /words", handlers.DeleteWordHandler)
	mux.HandleFunc("DELETE /words/{id}", handlers.DeleteOneWordHandler)
	return mux
}
