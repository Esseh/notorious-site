package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	PATH_NOTES_New      = "/new"
	PATH_NOTES_View     = "/view/:ID"
	PATH_NOTES_Editor   = "/edit/:ID"
	PATH_NOTES_Edit     = "/edit/"
)

func INIT_NOTES_HANDLERS(r *httprouter.Router) {
	r.GET(PATH_NOTES_New, NOTES_GET_New)
	r.POST(PATH_NOTES_New, NOTES_POST_New)
	r.GET(PATH_NOTES_View, NOTES_GET_View)
	r.GET(PATH_NOTES_Editor, NOTES_GET_Editor)
	r.POST(PATH_NOTES_Edit, NOTES_POST_Editor)
}


func NOTES_GET_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	_,validated := MustLogin(res, req); if !validated { return }
	ServeTemplateWithParams(res, "new-note", MakeHeader(res, req, false, true))
}

func NOTES_POST_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, validated := MustLogin(res,req); if !validated { return }
	ctx := NewContext(res,req)

	protected, boolConversionError := strconv.ParseBool(req.FormValue("protection"))
	if ErrorPage(ctx, res, nil, "Internal Server Error (1)", boolConversionError, http.StatusSeeOther) { return }
	
	NewContent := Content{
		Title:   req.FormValue("title"),
		Content: req.FormValue("note"),
	}

	key, err := retrievable.PlaceEntity(ctx, int64(0), &NewContent)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) { return }

	NewNote := Note{
		OwnerID:   int64(user.IntID),
		Protected: protected,
		ContentID: key.IntID(),
	}

	newkey, err := retrievable.PlaceEntity(ctx, int64(0), &NewNote)
	if ErrorPage(ctx, res, nil, "Internal Server Error (3)", err, http.StatusSeeOther) { return }
	http.Redirect(res, req, "/view/"+strconv.FormatInt(newkey.IntID(), 10), http.StatusSeeOther)
}

func NOTES_GET_View(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, validated := MustLogin(res,req);  if !validated { return }
	ctx := NewContext(res,req)

	NoteKey, intConversionError := strconv.ParseInt(params.ByName("ID"), 10, 64)
	if ErrorPage(ctx, res, nil, "Internal Server Error (1)", intConversionError, http.StatusSeeOther) { return }

	ViewNote := &Note{}
	ViewContent := &Content{}

	err := retrievable.GetEntity(ctx, NoteKey, ViewNote)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) { return }

	err = retrievable.GetEntity(ctx, ViewNote.ContentID, ViewContent)
	if ErrorPage(ctx, res, nil, "Internal Server Error (3)", err, http.StatusSeeOther) { return }

	owner, err := GetUserFromID(ctx, ViewNote.OwnerID)
	if ErrorPage(ctx, res, nil, "Internal Server Error (4)", err, http.StatusSeeOther) { return }

	NoteBody := template.HTML(EscapeString(ViewContent.Content))

	ServeTemplateWithParams(res, "viewNote", struct {
		HeaderData
		ErrorResponse, RedirectURL, Title, Notekey string
		Content                                    template.HTML
		User, Owner                                *User
	}{
		HeaderData:    *MakeHeader(res, req, false, true),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
		Title:         ViewContent.Title,
		Notekey:       params.ByName("ID"),
		Content:       NoteBody,
		User:          user,
		Owner:         owner,
	})

}









func NOTES_GET_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user, validated := MustLogin(res,req); if !validated { return }
	ctx := NewContext(res,req)
	ViewNote, ViewContent, err := GetNoteData(params.ByName("ID"), ctx)
	if ErrorPage(ctx, res, nil, "Internal Server Error (1)", err, http.StatusSeeOther) { return }

	validated := VerifyNotePermission(res, req, user, ViewNote); if !validated { return }

	Body := template.HTML(ViewContent.Content)
	ServeTemplateWithParams(res, "editnote", struct {
		HeaderData
		ErrorResponse, RedirectURL, Title, Notekey string
		Content                                    template.HTML
	}{
		HeaderData:    *MakeHeader(res, req, false, true),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
		Title:         ViewContent.Title,
		Notekey:       params.ByName("ID"),
		Content:       Body,
	})
}


















/// TODO: implement
func NOTES_POST_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	u, err := GetUserFromSession(req) // Check if a user is already logged in.
	//TODO DETERMInE IF THEY HAVE PERMISSION	// Should check in both GET and POST
	ctx := appengine.NewContext(req)
	if err != nil {
		http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
		return
	}

	data := req.FormValue("note")
	title := req.FormValue("title")
	notekey := req.FormValue("notekey")
	protection := req.FormValue("protection")
	log.Infof(ctx, "protections string is :", protection)
	protbool, err := strconv.ParseBool(protection)
	if ErrorPage(ctx, res, nil, "Internal Server Error (5)", err, http.StatusSeeOther) {
		return
	}

	Note := &Note{}

	intkey, err := strconv.ParseInt(notekey, 10, 64)

	err = retrievable.GetEntity(ctx, intkey, Note)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	// Permission Check, For collaberation it can also check against a collaborator container after the user check.
	// When setting for example, privacy setting might only be able to be set by the Owner so a separation is still needed.
	if Note.OwnerID != int64(u.IntID) && Note.Protected {
		// Soft rejection. can also be substituted for a http Not Allowed.
		http.Redirect(res, req, "/view/"+notekey, http.StatusSeeOther)
		return
	}

	Content := &Content{}

	err = retrievable.GetEntity(ctx, Note.ContentID, Content)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	tempcontent := EscapeString(data)

	Content.Content = tempcontent
	Content.Title = title
	if Note.OwnerID == int64(u.IntID) {
		Note.Protected = protbool
	}

	_, err = retrievable.PlaceEntity(ctx, intkey, Note)
	if ErrorPage(ctx, res, nil, "Internal Server Error (3)", err, http.StatusSeeOther) {
		return
	}

	_, err = retrievable.PlaceEntity(ctx, Note.ContentID, Content)
	if ErrorPage(ctx, res, nil, "Internal Server Error (4)", err, http.StatusSeeOther) {
		return
	}

	http.Redirect(res, req, "/view/"+notekey, http.StatusSeeOther)
}