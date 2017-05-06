package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Esseh/notorious-dev/AUTH"
	"github.com/Esseh/notorious-dev/BACKUP"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/NOTES"
	"github.com/Esseh/notorious-dev/PATHS"
	"github.com/Esseh/notorious-dev/RATINGS"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
)

func INIT_NOTES_HANDLERS(r *httprouter.Router) {
	r.GET(PATHS.NOTES_New, NOTES_GET_New)
	r.POST(PATHS.NOTES_New, NOTES_POST_New)
	r.GET(PATHS.NOTES_View, NOTES_GET_View)
	r.GET(PATHS.NOTES_Editor, NOTES_GET_Editor)
	r.POST(PATHS.NOTES_Edit, NOTES_POST_Editor)
	r.GET("/backup/:NoteID", NOTES_GET_Backups)
	r.POST("/note/api/getrating", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		fmt.Fprint(res, RATINGS.GetRating(CONTEXT.NewContext(res, req)))
	})
	r.POST("/note/api/setrating", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		fmt.Fprint(res, RATINGS.SetRating(CONTEXT.NewContext(res, req)))
	})
	r.POST("/note/api/copynote", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		fmt.Fprint(res, NOTES.API_SaveCopy(CONTEXT.NewContext(res, req)))
	})
}

func NOTES_GET_Backups(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	if !ctx.AssertLoggedInFailed() {
		// Good ID
		i, err := strconv.ParseInt(params.ByName("NoteID"), 10, 64)
		if err == nil {
			// Note Exists
			NoteData, _, _ := NOTES.GetExistingNote(ctx, params.ByName("NoteID"))
			err := retrievable.GetEntity(ctx, i, &NOTES.Note{})
			if err == nil {
				ref := AUTH.EmailReference{}
				if retrievable.GetEntity(ctx, req.FormValue("TargetEmail"), &ref) != nil {
					ref = AUTH.EmailReference{}
				}
				owner, _ := GetUserFromID(ctx, ref.UserID)
				ownerNote, _ := GetUserFromID(ctx, NoteData.OwnerID)
				b := BACKUP.RetrieveBackup(ctx, i, ref.UserID)
				CORE.ServeTemplateWithParams(res, "getbackups", struct {
					HeaderData CONTEXT.HeaderData
					Backup     NOTES.Content
					Backuphtml template.HTML
					Note       *NOTES.Note
					Owner      *USERS.User
					OwnerNote  *USERS.User
				}{*MakeHeader(ctx), b, template.HTML(b.Content), NoteData, owner, ownerNote})
			}
		}
	}
}

func NOTES_GET_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	if !ctx.AssertLoggedInFailed() {
		CORE.ServeTemplateWithParams(res, "new-note", struct{ HeaderData CONTEXT.HeaderData }{*MakeHeader(ctx)})
	}
}

func NOTES_POST_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	if !ctx.AssertLoggedInFailed() {
		publicedit, boolConversionError := strconv.ParseBool(req.FormValue("publicedit"))
		publicview, _ := strconv.ParseBool(req.FormValue("publicview"))
		if !ctx.ErrorPage("Internal Server Error (1)", boolConversionError, http.StatusSeeOther) {
			_, noteKey, err := NOTES.CreateNewNote(ctx,
				NOTES.Content{
					Title:   req.FormValue("title"),
					Content: req.FormValue("note"),
				},
				NOTES.Note{
					OwnerID:            int64(ctx.User.IntID),
					Collaborators:      FindCollaborators(ctx, req.FormValue("collaborators")),
					PublicallyEditable: publicedit,
					PublicallyViewable: publicview,
				},
			)
			if !ctx.ErrorPage("Internal Server Error (2)", err, http.StatusSeeOther) {
				ctx.Redirect("/view/" + strconv.FormatInt(noteKey.IntID(), 10))
			}
		}
	}
}

func NOTES_GET_View(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	if !ctx.AssertLoggedInFailed() {
		ViewNote, ViewContent, err := NOTES.GetExistingNote(ctx, params.ByName("ID"))
		if !ctx.ErrorPage("Internal Server Error (1)", err, http.StatusSeeOther) {
			if !NOTES.CanViewNote(ViewNote, ctx.User) {
				ctx.Redirect("/")
				return
			}
			owner, err := GetUserFromID(ctx, ViewNote.OwnerID)
			if !ctx.ErrorPage("Internal Server Error (2)", err, http.StatusSeeOther) {
				if !(NOTES.CanViewNote(ViewNote, ctx.User)) {
					ctx.Redirect("/")
				}
				NoteBody := template.HTML(CORE.EscapeString(ViewContent.Content))
				CORE.ServeTemplateWithParams(res, "viewNote", struct {
					HeaderData                                 CONTEXT.HeaderData
					ErrorResponse, RedirectURL, Title, Notekey string
					Content                                    template.HTML
					User, Owner                                *USERS.User
					NoteData                                   *NOTES.Note
				}{
					HeaderData:    *MakeHeader(ctx),
					RedirectURL:   req.FormValue("redirect"),
					ErrorResponse: req.FormValue("ErrorResponse"),
					Title:         ViewContent.Title,
					Notekey:       params.ByName("ID"),
					Content:       NoteBody,
					User:          ctx.User,
					Owner:         owner,
					NoteData:      ViewNote,
				})
			}
		}
	}
}

func NOTES_GET_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	if !ctx.AssertLoggedInFailed() {
		ViewNote, ViewContent, err := NOTES.GetExistingNote(ctx, params.ByName("ID"))
		if !ctx.ErrorPage("Internal Server Error (1)", err, http.StatusSeeOther) {
			owner, err := GetUserFromID(ctx, ViewNote.OwnerID)
			if !ctx.ErrorPage("Internal Server Error (2)", err, http.StatusSeeOther) {
				if !(NOTES.CanEditNote(ViewNote, ctx.User)) {
					ctx.Redirect("/")
				}
				Body := template.HTML(ViewContent.Content)
				CORE.ServeTemplateWithParams(res, "editnote", struct {
					HeaderData                                 CONTEXT.HeaderData
					ErrorResponse, RedirectURL, Title, Notekey string
					Content                                    template.HTML
					User, Owner                                *USERS.User
					NoteData                                   *NOTES.Note
				}{
					HeaderData:    *MakeHeader(ctx),
					RedirectURL:   req.FormValue("redirect"),
					ErrorResponse: req.FormValue("ErrorResponse"),
					Title:         ViewContent.Title,
					Notekey:       params.ByName("ID"),
					Content:       Body,
					NoteData:      ViewNote,
					User:          ctx.User,
					Owner:         owner,
				})
			}
		}
	}
}

func NOTES_POST_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := CONTEXT.NewContext(res, req)
	ctx.Redirect("/view/" + req.FormValue("notekey"))
	if !ctx.AssertLoggedInFailed() {
		edit, boolConversionError := strconv.ParseBool(req.FormValue("publicedit"))
		view, _ := strconv.ParseBool(req.FormValue("publicview"))
		if !ctx.ErrorPage("Internal Server Error (1)", boolConversionError, http.StatusSeeOther) {
			content := NOTES.Content{
				Content: CORE.EscapeString(req.FormValue("note")),
				Title:   req.FormValue("title"),
			}
			err := NOTES.UpdateNoteContent(ctx, req.FormValue("notekey"),
				content,
				NOTES.Note{
					Collaborators:      FindCollaborators(ctx, req.FormValue("collaborators")),
					PublicallyEditable: edit,
					PublicallyViewable: view,
				},
			)
			noteid, _ := strconv.ParseInt(req.FormValue("notekey"), 10, 64)
			n := NOTES.Note{}
			retrievable.GetEntity(ctx, noteid, &n)
			isCollaborator := false
			for _, v := range n.Collaborators {
				if int64(ctx.User.IntID) == n.OwnerID || int64(ctx.User.IntID) == v {
					isCollaborator = true
					break
				}
			}
			if isCollaborator {
				b := BACKUP.Backup(content)
				BACKUP.UpdateBackup(ctx, noteid, int64(ctx.User.IntID), &b)
			}
			if !ctx.ErrorPage("Internal Server Error (2)", err, http.StatusSeeOther) {
				ctx.Redirect("/view/" + req.FormValue("notekey"))
			}
		}
	}
}
