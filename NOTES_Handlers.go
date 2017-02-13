package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
	ctx := NewContext(res,req)
	if ctx.AssertLoggedInFailed() { return }
	ServeTemplateWithParams(res, "new-note", MakeHeader(ctx))
}

func NOTES_POST_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if ctx.AssertLoggedInFailed() { return }
	
	protected, boolConversionError := strconv.ParseBool(req.FormValue("protection"))
	if ErrorPage(ctx, "Internal Server Error (1)", boolConversionError, http.StatusSeeOther) { return }

	
	_, noteKey, err := CreateNewNote(ctx,
		Content{
			Title:   req.FormValue("title"),
			Content: req.FormValue("note"),
		},
		Note{
			OwnerID:   int64(ctx.user.IntID),
			Protected: protected,
		},
	)
	if ErrorPage(ctx, "Internal Server Error (2)", err, http.StatusSeeOther) { return }
	
	ctx.Redirect("/view/"+strconv.FormatInt(noteKey.IntID(), 10))
}

func NOTES_GET_View(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if ctx.AssertLoggedInFailed() { return }
	
	ViewNote, ViewContent, err := GetExistingNote(ctx,params.ByName("ID"))
	if ErrorPage(ctx, "Internal Server Error (1)", err, http.StatusSeeOther) { return }

	owner, err := GetUserFromID(ctx, ViewNote.OwnerID)
	if ErrorPage(ctx, "Internal Server Error (2)", err, http.StatusSeeOther) { return }

	NoteBody := template.HTML(EscapeString(ViewContent.Content))

	ServeTemplateWithParams(res, "viewNote", struct {
		HeaderData
		ErrorResponse, RedirectURL, Title, Notekey string
		Content                                    template.HTML
		User, Owner                                *User
	}{
		HeaderData:    *MakeHeader(ctx),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
		Title:         ViewContent.Title,
		Notekey:       params.ByName("ID"),
		Content:       NoteBody,
		User:          ctx.user,
		Owner:         owner,
	})

}

func NOTES_GET_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if ctx.AssertLoggedInFailed() { return }
	ViewNote, ViewContent, err := GetExistingNote(ctx,params.ByName("ID"))
	if ErrorPage(ctx, "Internal Server Error (1)", err, http.StatusSeeOther) { return }

	validated := VerifyNotePermission(ctx, ViewNote); if !validated { return }

	Body := template.HTML(ViewContent.Content)
	ServeTemplateWithParams(res, "editnote", struct {
		HeaderData
		ErrorResponse, RedirectURL, Title, Notekey string
		Content                                    template.HTML
	}{
		HeaderData:    *MakeHeader(ctx),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
		Title:         ViewContent.Title,
		Notekey:       params.ByName("ID"),
		Content:       Body,
	})
}

func NOTES_POST_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := NewContext(res,req)
	if ctx.AssertLoggedInFailed() { return }
	
	protbool, boolConversionError := strconv.ParseBool(req.FormValue("protection"))
	if ErrorPage(ctx, "Internal Server Error (1)", boolConversionError, http.StatusSeeOther) { return }

	err := UpdateNoteContent(ctx,req.FormValue("notekey"),
		Content{
			Content: EscapeString(req.FormValue("note")),
			Title: req.FormValue("title"),
		},
		Note{
			Protected: protbool,
		},
	)
	if ErrorPage(ctx, "Internal Server Error (2)", err, http.StatusSeeOther) { return }	
	ctx.Redirect("/view/"+req.FormValue("notekey"))
}