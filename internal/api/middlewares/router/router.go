package router

import "net/http"

func MainRouter() *http.ServeMux {

	//subroutes
	eRouter := execsRouter()
	dRouter := dictionaryRouter()

	dRouter.Handle("/", eRouter)
	return dRouter

}
