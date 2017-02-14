package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

const DomainPath = "http://localhost:8080/"

// Multiplexer Function for CORE
func Handle_CORE(r *httprouter.Router) {
	r.GET("/", index)
}

// Serves the index page.
func index(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ctx := NewContext(res,req)
	ServeTemplateWithParams(res, "index", struct { HeaderData }{ *MakeHeader(ctx), })
}

func init() {
	router := httprouter.New()
	InitializeHandlers(router)
	http.Handle("/", router)
}

func InitializeHandlers(router *httprouter.Router) {
	Handle_CORE(router)
	INIT_AUTH_HANDLERS(router)
	INIT_OAUTH_Handlers(router)
	INIT_USERS_HANDLERS(router)
	INIT_NOTES_HANDLERS(router)
}