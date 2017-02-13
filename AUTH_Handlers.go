package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

const (
	PATH_AUTH_Login          = "/login"
	PATH_AUTH_Logout         = "/logout"
	PATH_AUTH_Register       = "/register"
)

func INIT_AUTH_HANDLERS(r *httprouter.Router) {
	r.GET(PATH_AUTH_Logout, AUTH_GET_Logout)                   
	r.GET(PATH_AUTH_Login, AUTH_GET_Login)                     
	r.POST(PATH_AUTH_Login, AUTH_POST_Login)                   
	r.GET(PATH_AUTH_Register, AUTH_GET_Register)               
	r.POST(PATH_AUTH_Register, AUTH_POST_Register)             
}

func AUTH_GET_Login(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	ServeTemplateWithParams(res, "login", struct {
		HeaderData
		ErrorResponse, RedirectURL string
	}{
		HeaderData:    *MakeHeader(ctx),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
	})
}

func AUTH_POST_Login(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if !ValidLogin(req.FormValue("email"), req.FormValue("password")) {
		BackWithError(ctx, ErrInvalidLogin, "Invalid Login Information")
	} else {
		response, err := LoginToWebsite(ctx,req.FormValue("email"), req.FormValue("password"))
		if !BackWithError(ctx, err, response) { ctx.Redirect("/"+req.FormValue("redirect")) }
	}
}

func AUTH_GET_Logout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	response,err := LogoutFromWebsite(ctx)	
	if !ErrorPage(ctx, response, err, http.StatusBadRequest) {  ctx.Redirect("/"+req.FormValue("redirect")) }
}

func AUTH_GET_Register(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	ServeTemplateWithParams(res, "register", struct {
		HeaderData
		ErrorResponse, RedirectURL string
	}{
		HeaderData:    *MakeHeader(ctx),
		ErrorResponse: req.FormValue("ErrorResponse"),
		RedirectURL:   req.FormValue("redirect"),
	})
}
func AUTH_POST_Register(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	response, err := RegisterNewUser(
		ctx,
		req.FormValue("email"),
		req.FormValue("password"),
		req.FormValue("cpassword"),
		req.FormValue("first"),
		req.FormValue("last"),
	)
	if !BackWithError(ctx, err, response) { AUTH_POST_Login(res, req, params) }
}