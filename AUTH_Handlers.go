package main

import (
	"errors"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/Esseh/notorious-dev/AUTH"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/PATHS"

)


func INIT_AUTH_HANDLERS(r *httprouter.Router) {
	r.GET(PATHS.AUTH_Logout, AUTH_GET_Logout)                   
	r.GET(PATHS.AUTH_Login, AUTH_GET_Login)                     
	r.POST(PATHS.AUTH_Login, AUTH_POST_Login)                   
	r.GET(PATHS.AUTH_Register, AUTH_GET_Register)               
	r.POST(PATHS.AUTH_Register, AUTH_POST_Register)             
}

func AUTH_GET_Login(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	CORE.ServeTemplateWithParams(res, "login", struct {
		HeaderData	CONTEXT.HeaderData
		ErrorResponse, RedirectURL string
	}{
		HeaderData:    *MakeHeader(ctx),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
	})
}

func AUTH_POST_Login(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	if !CORE.ValidLogin(req.FormValue("email"), req.FormValue("password")) {
		ctx.BackWithError(errors.New("Not Logged In"), "Invalid Login Information")
	} else {
		response, err := AUTH.LoginToWebsite(ctx,req.FormValue("email"), req.FormValue("password"))
		if !ctx.BackWithError(err, response) { ctx.Redirect("/"+req.FormValue("redirect")) }
	}
}

func AUTH_GET_Logout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	response,err := AUTH.LogoutFromWebsite(ctx)	
	if !ctx.ErrorPage(response, err, http.StatusBadRequest) {  ctx.Redirect("/"+req.FormValue("redirect")) }
}

func AUTH_GET_Register(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	CORE.ServeTemplateWithParams(res, "register", struct {
		HeaderData CONTEXT.HeaderData
		ErrorResponse, RedirectURL string
	}{
		HeaderData:    *MakeHeader(ctx),
		ErrorResponse: req.FormValue("ErrorResponse"),
		RedirectURL:   req.FormValue("redirect"),
	})
}
func AUTH_POST_Register(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	response, err := AUTH.RegisterNewUser(
		ctx,
		req.FormValue("email"),
		req.FormValue("password"),
		req.FormValue("cpassword"),
		req.FormValue("first"),
		req.FormValue("last"),
	)
	if !ctx.BackWithError(err, response) { AUTH_POST_Login(res, req, params) }
}