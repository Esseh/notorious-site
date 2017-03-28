// The USERS Module, Deals with the User interfacing with themselves.
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/NOTES"
	"github.com/Esseh/notorious-dev/PATHS"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
)

func INIT_USERS_HANDLERS(r *httprouter.Router) {
	r.GET(PATHS.USERS_ProfileEdit, USERS_GET_ProfileEdit)
	r.POST(PATHS.USERS_ProfileEdit, USERS_POST_ProfileEdit)
	r.POST(PATHS.USERS_ProfileEditAvatar, USERS_POST_ProfileEditAvatar)
	r.GET(PATHS.USERS_ProfileView, USERS_GET_ProfileView)
}

func USERS_GET_ProfileView(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	id, convErr := strconv.ParseInt(params.ByName("ID"), 10, 64)
	if !ctx.ErrorPage("Invalid ID", convErr, http.StatusBadRequest) {
		ci, getErr := GetUserFromID(ctx, id)
		if !ctx.ErrorPage("Not a valid user ID", getErr, http.StatusNotFound) {
			notes, err := NOTES.GetAllNotes(ctx, id)
			if !ctx.ErrorPage("Internal Server Error", err, http.StatusSeeOther) {
				avatarMod := GetMod(id)
				screen := struct {
					CONTEXT.HeaderData
					Data      *USERS.User
					AllNotes  []NOTES.NoteOutput
					AvatarMod int64
					Root	  int64
				}{
					*MakeHeader(ctx),
					ci,
					notes,
					avatarMod,
					id,
				}
				CORE.ServeTemplateWithParams(res, "user-profile", screen)
			}
		}
	}
}

func USERS_GET_ProfileEdit(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
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
	ctx := CONTEXT.NewContext(res, req)
	user, _ := USERS.GetUserFromSession(ctx, req)
	user.First = req.FormValue("first")
	user.Last = req.FormValue("last")
	user.Bio = req.FormValue("bio")
	_, err := retrievable.PlaceEntity(ctx, user.IntID, user)
	if !ctx.ErrorPage("server error placing key", err, http.StatusBadRequest) {
		ctx.Redirect("/profile/" + strconv.FormatInt(int64(user.IntID), 10))
	}
}

func USERS_POST_ProfileEditAvatar(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	user, _ := USERS.GetUserFromSession(ctx, req)
	avatarImg, hdr, err := req.FormFile("avatar")
	if err != nil {
		return
	}
	defer avatarImg.Close()
	posX, _ := strconv.Atoi(req.FormValue("posx"))
	posY, _ := strconv.Atoi(req.FormValue("posy"))
	cropWidth, err := strconv.Atoi(req.FormValue("cropwidth"))
	if err != nil {
		cropWidth = 500
	}
	cropHeight, err := strconv.Atoi(req.FormValue("cropheight"))
	if err != nil {
		cropHeight = 500
	}
	rotate, _ := strconv.Atoi(req.FormValue("degrees"))
	cb := CropBounds{
		X:         posX,
		Y:         posY,
		W:         cropWidth,
		H:         cropHeight,
		RotateDeg: rotate,
	}
	err2 := uploadImage(ctx, int64(user.IntID), hdr, &cb, avatarImg)
	if err2 != nil {
		fmt.Fprint(res, err2)
	} else {
		user.Avatar = true
		_, err = retrievable.PlaceEntity(ctx, user.IntID, user)
		if !ctx.ErrorPage("server error placing key", err, http.StatusBadRequest) {
			http.Redirect(res, req, "/profile/"+strconv.FormatInt(int64(user.IntID), 10), http.StatusSeeOther)
		}
	}
}
