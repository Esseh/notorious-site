package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

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