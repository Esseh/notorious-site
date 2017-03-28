package main

import (
	"net/http"

	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/julienschmidt/httprouter"
)

// Multiplexer Function for CORE
func Handle_CORE(r *httprouter.Router) {
	r.GET("/", index)
}

// Serves the index page.
func index(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	CORE.ServeTemplateWithParams(res, "index", struct{ CONTEXT.HeaderData }{*MakeHeader(ctx)})
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
