// The USERS Module, Deals with the User interfacing with themselves.
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
)

func INIT_USERS_HANDLERS(r *httprouter.Router) {
	r.GET(PATH_USERS_ProfileEdit, USERS_GET_ProfileEdit)
	r.POST(PATH_USERS_ProfileEdit, USERS_POST_ProfileEdit)
	r.POST(PATH_USERS_ProfileEditAvatar, USERS_POST_ProfileEditAvatar)
	r.GET(PATH_USERS_ProfileView, USERS_GET_ProfileView)
}

const (
	PATH_USERS_ProfileEdit       = "/editprofile"
	PATH_USERS_ProfileEditAvatar = "/editprofileavatar"
	PATH_USERS_ProfileView       = "/profile/:ID"
)

//===========================================================================
// Profile
//===========================================================================
func USERS_GET_ProfileEdit(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if ctx.AssertLoggedInFailed() { return }
	err := ServeTemplateWithParams(res, "profile-settings", struct {
		HeaderData
		ErrorResponseProfile string
		User                 *User
	}{
		*MakeHeader(ctx),
		req.FormValue("ErrorResponseProfile"),
		ctx.user,
	})
	if err != nil {
		fmt.Fprint(res, err)
	}
}

// TODO: Implement
func USERS_POST_ProfileEdit(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	u, _ := AUTH_GetUserFromSession(req)
	u.First = req.FormValue("first")
	u.Last = req.FormValue("last")
	u.Bio = req.FormValue("bio")
	ctx := NewContext(res,req)
	_, err := retrievable.PlaceEntity(ctx, u.IntID, u)
	if ErrorPage(ctx, "server error placing key", err, http.StatusBadRequest) {
		return
	}

	http.Redirect(res, req, "/profile/"+strconv.FormatInt(int64(u.IntID), 10), http.StatusSeeOther)
}

// TODO: Implement
func USERS_POST_ProfileEditAvatar(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	u, _ := AUTH_GetUserFromSession(req)
	ctx := NewContext(res,req)

	rdr, hdr, err := req.FormFile("avatar")
	if ErrorPage(ctx, "upload image thingy", err, http.StatusBadRequest) {
		return
	}
	defer rdr.Close()
	u.Avatar = true
	err2 := UploadAvatar(ctx, int64(u.IntID), hdr, rdr)

	if err2 != nil {
		fmt.Fprint(res, err2)
	}
	_, err = retrievable.PlaceEntity(ctx, u.IntID, u)
	if ErrorPage(ctx, "server error placing key", err, http.StatusBadRequest) {
		return
	}
	http.Redirect(res, req, "/profile/"+strconv.FormatInt(int64(u.IntID), 10), http.StatusSeeOther)
}

//===========================================================================
// Profile View
//===========================================================================
func USERS_GET_ProfileView(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	id, convErr := strconv.ParseInt(params.ByName("ID"), 10, 64)
	if ErrorPage(ctx, "Invalid ID", convErr, http.StatusBadRequest) {
		return
	}
	ci, getErr := AUTH_GetUserFromID(ctx, id)
	if ErrorPage(ctx, "Not a valid user ID", getErr, http.StatusNotFound) {
		return
	}
	notes, err := GetAllNotes(ctx, id)
	if ErrorPage(ctx, "Internal Server Error", err, http.StatusSeeOther) {
		return
	}
	screen := struct {
		HeaderData
		Data     *User
		AllNotes []NoteOutput
	}{
		*MakeHeader(ctx),
		ci,
		notes,
	}
	ServeTemplateWithParams(res, "user-profile", screen)
}
