// The USERS Module, Deals with the User interfacing with themselves.
package main

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/NOTES"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/PATHS"
	"github.com/julienschmidt/httprouter"
)

func INIT_USERS_HANDLERS(r *httprouter.Router) {
	r.GET(PATHS.USERS_ProfileEdit, USERS_GET_ProfileEdit)
	r.POST(PATHS.USERS_ProfileEdit, USERS_POST_ProfileEdit)
	r.POST(PATHS.USERS_ProfileEditAvatar, USERS_POST_ProfileEditAvatar)
	r.GET(PATHS.USERS_ProfileView, USERS_GET_ProfileView)
}

func USERS_GET_ProfileView(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	id, convErr := strconv.ParseInt(params.ByName("ID"), 10, 64)
	if !ctx.ErrorPage("Invalid ID", convErr, http.StatusBadRequest) {
		ci, getErr := GetUserFromID(ctx, id)
		if !ctx.ErrorPage("Not a valid user ID", getErr, http.StatusNotFound) {
			notes, err := NOTES.GetAllNotes(ctx, id)
			if !ctx.ErrorPage("Internal Server Error", err, http.StatusSeeOther) {
				screen := struct {
					CONTEXT.HeaderData
					Data     *USERS.User
					AllNotes []NOTES.NoteOutput
				}{
					*MakeHeader(ctx),
					ci,
					notes,
				}
				CORE.ServeTemplateWithParams(res, "user-profile", screen)
			}
		}
	}
}


func USERS_GET_ProfileEdit(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	if !ctx.AssertLoggedInFailed() {
		CORE.ServeTemplateWithParams(res, "profile-settings", struct {
			CONTEXT.HeaderData
			ErrorResponseProfile string
			User                 *USERS.User
		}{
			*MakeHeader(ctx),
			req.FormValue("ErrorResponseProfile"),
			ctx.User,
		})
	}
}

func USERS_POST_ProfileEdit(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	user, _ := USERS.GetUserFromSession(ctx,req)
	user.First = req.FormValue("first")
	user.Last = req.FormValue("last")
	user.Bio = req.FormValue("bio")
	_, err := retrievable.PlaceEntity(ctx, user.IntID, user)
	if !ctx.ErrorPage("server error placing key", err, http.StatusBadRequest) {
		ctx.Redirect("/profile/"+strconv.FormatInt(int64(user.IntID), 10))
	}
}

func USERS_POST_ProfileEditAvatar(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	user, _ := USERS.GetUserFromSession(ctx,req)
	rdr, hdr, err := req.FormFile("avatar")
	if !ctx.ErrorPage("upload image thingy", err, http.StatusBadRequest) {
		defer rdr.Close()
		user.Avatar = true
		err2 := USERS.UploadAvatar(ctx, int64(user.IntID), hdr, rdr)
		if err2 != nil { fmt.Fprint(res, err2) } else {
			_, err = retrievable.PlaceEntity(ctx, user.IntID, user)
			if !ctx.ErrorPage("server error placing key", err, http.StatusBadRequest) {
				http.Redirect(res, req, "/profile/"+strconv.FormatInt(int64(user.IntID), 10), http.StatusSeeOther)
			}
		}
	}
}