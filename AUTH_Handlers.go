// Handlers dealing with authentication
package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
)

const (
	PATH_AUTH_Login          = "/login"
	PATH_AUTH_Logout         = "/logout"
	PATH_AUTH_Register       = "/register"
)

func INIT_AUTH_HANDLERS(r *httprouter.Router) {
	r.GET(PATH_AUTH_Logout, AUTH_GET_Logout)                   // PATH_AUTH_Logout 				= "/logout"
	r.GET(PATH_AUTH_Login, AUTH_GET_Login)                     // PATH_AUTH_Login 				= "/login"
	r.POST(PATH_AUTH_Login, AUTH_POST_Login)                   //
	r.GET(PATH_AUTH_Register, AUTH_GET_Register)               // PATH_AUTH_Register 			= "/register"
	r.POST(PATH_AUTH_Register, AUTH_POST_Register)             //
}

//=========================================================================================
// Login
//=========================================================================================
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
	username := strings.ToLower(req.FormValue("email"))
	password := req.FormValue("password")
	redirect := req.FormValue("redirect")

	if username == "" || password == "" { // Check incoming information for Trivial False case.
		v := url.Values{}
		v.Add("redirect", redirect)
		v.Add("ErrorResponse", "Fields Cannot Be Empty")
		http.Redirect(res, req, PATH_AUTH_Login+"?"+v.Encode(), http.StatusSeeOther)
		return
	}

	userID, err := GetUserIDFromLogin(ctx, username, password)
	if BackWithError(res, req, err, "Login Information Is Incorrect") {
		return
	}

	sessionID, err := CreateSessionID(ctx, req, userID)
	if BackWithError(res, req, err, "Login error, try again later.") {
		return
	}

	err = MakeCookie(res, "session", strconv.FormatInt(sessionID, 10))
	if BackWithError(res, req, err, "Login error, try again later.") {
		return
	}
	http.Redirect(res, req, "/"+redirect, http.StatusSeeOther)
}

//=========================================================================================
//Logout
//=========================================================================================
func AUTH_GET_Logout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	sessionIDStr, err := GetCookieValue(req, "session")
	if ErrorPage(ctx, "Must be logged in", err, http.StatusBadRequest) {
		return
	}

	sessionVal, err := strconv.ParseInt(sessionIDStr, 10, 0)
	if ErrorPage(ctx, "Bad cookie value", err, http.StatusBadRequest) {
		return
	}

	err = retrievable.DeleteEntity(ctx, (&Session{}).Key(ctx, sessionVal))
	if ErrorPage(ctx, "No such session found!", err, 500) {
		return
	}

	DeleteCookie(res, "session")
	http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
}

//=========================================================================================
//Register
//=========================================================================================
func AUTH_GET_Register(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if ctx.AssertLoggedInFailed() { return }
	ServeTemplateWithParams(res, "register", struct {
		HeaderData
		BusinessKey, ErrorResponse, RedirectURL string
	}{
		HeaderData:    *MakeHeader(ctx),
		ErrorResponse: req.FormValue("ErrorResponse"),
		BusinessKey:   req.FormValue("BusinessKey"),
		RedirectURL:   req.FormValue("redirect"),
	})
}
func AUTH_POST_Register(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	nu := &User{ // Make the New User
		Email:    strings.ToLower(req.FormValue("email")),
		First:    req.FormValue("first"),
		Last:     req.FormValue("last"),
	}

	password := req.FormValue("password")
	confirmPassword := req.FormValue("cpassword")

	// Check for trivially false.
	if "" == nu.Email || "" == nu.First || "" == nu.Last || "" == password || "" == confirmPassword || password != confirmPassword {
		BackWithError(res, req, ErrEmptyField, "Field Empty or password mismatch")
		return
	}

	nu, err := CreateUserFromLogin(ctx, nu.Email, password, nu)
	if BackWithError(res, req, err, "Username Taken") {
		return
	}

	AUTH_POST_Login(res, req, params)
}