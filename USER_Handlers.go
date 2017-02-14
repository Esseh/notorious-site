// The USERS Module, Deals with the User interfacing with themselves.
package main

import (
	"net/http"
	"strconv"
	"fmt"
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

func USERS_GET_ProfileView(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	id, convErr := strconv.ParseInt(params.ByName("ID"), 10, 64)
	if !ERROR_Page(ctx, "Invalid ID", convErr, http.StatusBadRequest) {
		ci, getErr := AUTH_GetUserFromID(ctx, id)
		if !ERROR_Page(ctx, "Not a valid user ID", getErr, http.StatusNotFound) {
			notes, err := NOTES_GetAllNotes(ctx, id)
			if !ERROR_Page(ctx, "Internal Server Error", err, http.StatusSeeOther) {
				screen := struct {
					HeaderData
					Data     *USER_User
					AllNotes []NOTES_NoteOutput
				}{
					*MakeHeader(ctx),
					ci,
					notes,
				}
				ServeTemplateWithParams(res, "user-profile", screen)
			}
		}
	}
}


func USERS_GET_ProfileEdit(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if !ctx.AssertLoggedInFailed() {
		ServeTemplateWithParams(res, "profile-settings", struct {
			HeaderData
			ErrorResponseProfile string
			User                 *USER_User
		}{
			*MakeHeader(ctx),
			req.FormValue("ErrorResponseProfile"),
			ctx.user,
		})
	}
}

func USERS_POST_ProfileEdit(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, _ := AUTH_GetUserFromSession(req)
	user.First = req.FormValue("first")
	user.Last = req.FormValue("last")
	user.Bio = req.FormValue("bio")
	ctx := NewContext(res,req)
	_, err := retrievable.PlaceEntity(ctx, user.IntID, user)
	if !ERROR_Page(ctx, "server error placing key", err, http.StatusBadRequest) {
		ctx.Redirect("/profile/"+strconv.FormatInt(int64(user.IntID), 10))
	}
}

func USERS_POST_ProfileEditAvatar(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, _ := AUTH_GetUserFromSession(req)
	ctx := NewContext(res,req)
	rdr, hdr, err := req.FormFile("avatar")
	if !ERROR_Page(ctx, "upload image thingy", err, http.StatusBadRequest) {
		defer rdr.Close()
		user.Avatar = true
		err2 := USER_UploadAvatar(ctx, int64(user.IntID), hdr, rdr)
		if err2 != nil { fmt.Fprint(res, err2) } else {
			_, err = retrievable.PlaceEntity(ctx, user.IntID, user)
			if !ERROR_Page(ctx, "server error placing key", err, http.StatusBadRequest) {
				http.Redirect(res, req, "/profile/"+strconv.FormatInt(int64(user.IntID), 10), http.StatusSeeOther)
			}
		}
	}
}