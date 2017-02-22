package main

import (
	"html/template"
	"net/http"
	"strconv"
	"github.com/Esseh/notorious-dev/PATHS"
	"github.com/julienschmidt/httprouter"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/NOTES"
)

func INIT_NOTES_HANDLERS(r *httprouter.Router) {
	r.GET(PATHS.NOTES_New, NOTES_GET_New)
	r.POST(PATHS.NOTES_New, NOTES_POST_New)
	r.GET(PATHS.NOTES_View, NOTES_GET_View)
	r.GET(PATHS.NOTES_Editor, NOTES_GET_Editor)
	r.POST(PATHS.NOTES_Edit, NOTES_POST_Editor)
}


func NOTES_GET_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	if !ctx.AssertLoggedInFailed() { CORE.ServeTemplateWithParams(res, "new-note", struct{HeaderData CONTEXT.HeaderData}{*MakeHeader(ctx)}) }
}

func NOTES_POST_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	if !ctx.AssertLoggedInFailed() {
		publicallyEditable, boolConversionError := strconv.ParseBool(req.FormValue("publicallyeditable"))
		if !ctx.ErrorPage("Internal Server Error (1)", boolConversionError, http.StatusSeeOther) {		
			publicallyViewable, boolConversionError := strconv.ParseBool(req.FormValue("publicallyeditable"))
			if !ctx.ErrorPage("Internal Server Error (3)", boolConversionError, http.StatusSeeOther) {	
				_, noteKey, err := NOTES.CreateNewNote(ctx,
					NOTES.Content{
						Title:   req.FormValue("title"),
						Content: req.FormValue("note"),
					},
					NOTES.Note{
						OwnerID:   int64(ctx.User.IntID),
						PublicallyViewable: publicallyViewable,
						PublicallyEditable: publicallyEditable,		
					},
				)
				if !ctx.ErrorPage("Internal Server Error (2)", err, http.StatusSeeOther) {
					ctx.Redirect("/view/"+strconv.FormatInt(noteKey.IntID(), 10))
				}
			}
		}
	}
}

func NOTES_GET_View(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	if !ctx.AssertLoggedInFailed() {
		ViewNote, ViewContent, err := NOTES.GetExistingNote(ctx,params.ByName("ID"))
		if !ctx.ErrorPage("Internal Server Error (1)", err, http.StatusSeeOther) {
			if !NOTES.CanViewNote(ViewNote,ctx.User){ ctx.Redirect("/"); return }
			owner, err := GetUserFromID(ctx, ViewNote.OwnerID)
			if !ctx.ErrorPage("Internal Server Error (2)", err, http.StatusSeeOther) {
				NoteBody := template.HTML(CORE.EscapeString(ViewContent.Content))
				CORE.ServeTemplateWithParams(res, "viewNote", struct {
					HeaderData CONTEXT.HeaderData
					ErrorResponse, RedirectURL, Title, Notekey string
					Content                                    template.HTML
					User, Owner                                *USERS.User
					NoteData								   *NOTES.Note
				}{
					HeaderData:    *MakeHeader(ctx),
					RedirectURL:   req.FormValue("redirect"),
					ErrorResponse: req.FormValue("ErrorResponse"),
					Title:         ViewContent.Title,
					Notekey:       params.ByName("ID"),
					Content:       NoteBody,
					User:          ctx.User,
					Owner:         owner,
					NoteData:	   ViewNote,
				})
			}
		}
	}
}

func NOTES_GET_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	if !ctx.AssertLoggedInFailed() { 
		ViewNote, ViewContent, err := NOTES.GetExistingNote(ctx,params.ByName("ID"))
		if !ctx.ErrorPage("Internal Server Error (1)", err, http.StatusSeeOther) {
			validated := NOTES.CanEditNote(ViewNote,ctx.User)
			if validated {
				Body := template.HTML(ViewContent.Content)
				CORE.ServeTemplateWithParams(res, "editnote", struct {
					HeaderData CONTEXT.HeaderData
					ErrorResponse, RedirectURL, Title, Notekey string
					Content                                    template.HTML
					NoteData								   *NOTES.Note
				}{
					HeaderData:    *MakeHeader(ctx),
					RedirectURL:   req.FormValue("redirect"),
					ErrorResponse: req.FormValue("ErrorResponse"),
					Title:         ViewContent.Title,
					Notekey:       params.ByName("ID"),
					Content:       Body,
					NoteData:	   ViewNote,
				})
			}
		}
	}
}

func NOTES_POST_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res,req)
	if !ctx.AssertLoggedInFailed() {
		publicallyViewable, boolConversionError := strconv.ParseBool(req.FormValue("publiclyviewable"))
		if !ctx.ErrorPage("Internal Server Error (1)", boolConversionError, http.StatusSeeOther) {
			publicallyEditable, boolConversionError := strconv.ParseBool(req.FormValue("publiclyeditable"))
			if !ctx.ErrorPage("Internal Server Error (3)", boolConversionError, http.StatusSeeOther) {
				err := NOTES.UpdateNoteContent(ctx,req.FormValue("notekey"),
					NOTES.Content{
						Content: CORE.EscapeString(req.FormValue("note")),
						Title: req.FormValue("title"),
					},
					NOTES.Note{
						PublicallyEditable: publicallyEditable,
						PublicallyViewable: publicallyViewable,
					},
				)
				if !ctx.ErrorPage("Internal Server Error (2)", err, http.StatusSeeOther) { 
					ctx.Redirect("/view/"+req.FormValue("notekey"))
				}
			}
		}
	}
}