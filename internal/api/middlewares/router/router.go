package router

import "net/http"

//TODO WATCH THIS MOUNTING

// main router sets up the HTTP router by chaning sub routes for auth and dict APIS
func MainRouter() *http.ServeMux {

	//subroutes
	eRouter := execsRouter()
	dRouter := dictionaryRouter()

	dRouter.Handle("/", eRouter)
	return dRouter

}
