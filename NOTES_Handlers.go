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
	
	PATH_NOTES_Reference_Test = "/references/test/:USER"
)

func INIT_NOTES_HANDLERS(r *httprouter.Router) {
	r.GET(PATH_NOTES_New, NOTES_GET_New)
	r.POST(PATH_NOTES_New, NOTES_POST_New)
	r.GET(PATH_NOTES_View, NOTES_GET_View)
	r.GET(PATH_NOTES_Editor, NOTES_GET_Editor)
	r.POST(PATH_NOTES_Edit, NOTES_POST_Editor)
	
	// Test For Matthew
	r.GET(PATH_NOTES_Reference_Test, NOTES_GET_Reference_TestBed)
}

func NOTES_GET_Reference_TestBed(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	/*
		HeaderData
			ShowLogin    bool
			ShowRegister bool
			Ctx          context.Context
			User         *User
			CurrentPath  string
	*/
	/*
		User
			First, Last       string
			Email             string
			Avatar            bool `datastore:",noindex"`
			Bio               string
			retrievable.IntID `datastore:"-" json:"-"`
	*/
	/*
		Note
			OwnerID int64 // owner of the note, can set permissions
			//Collab []int64			// any people collabing, stretch goal
			Protected bool  // if true then it is not publically editable.
			ContentID int64 // Reference to the content of the note.

		Content
			Title, Content string // Content can be raw html or markdown	
	*/
	obj := struct {
		HeaderData
		MainNote		 Note
		ReferenceNotes []Note
	}{
		*MakeHeader(res, req, false, true),
		Note{5,false,0},
		[]Note{
			Note{1,true,0},		
			Note{2,false,0},		
			Note{2,true,0},		
			Note{2,false,0},		
		},
	}
	// test with existing user
	if(params.ByName("ID") == "1"){
		obj.HeaderData.User = &User{"Bob","Dole","mail@mail.com",false,"Just a man looking for love.",retrievable.IntID(1)}
	} else {
	// test without existing user
		obj.HeaderData.User = &User{}
	}
	ServeTemplateWithParams(res, "note-test", obj)
}

func NOTES_GET_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if MustLogin(res, req) {
		return
	}

	ServeTemplateWithParams(res, "newnote", struct {
		HeaderData
		ErrorResponse, RedirectURL string
	}{
		HeaderData:    *MakeHeader(res, req, false, true),
		RedirectURL:   req.FormValue("redirect"),
		ErrorResponse: req.FormValue("ErrorResponse"),
	})
}

func NOTES_POST_New(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	u, err := GetUserFromSession(req) // Check if a user is already logged in.
	ctx := appengine.NewContext(req)

	if err != nil {
		http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
		return
	}

	data := req.FormValue("note")
	title := req.FormValue("title")
	protected, boolerr := strconv.ParseBool(req.FormValue("protection"))
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", boolerr, http.StatusSeeOther) {
		return
	}

	NewContent := Content{
		Title:   title,
		Content: data,
	}

	key, err := retrievable.PlaceEntity(ctx, int64(0), &NewContent)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	NewNote := Note{
		OwnerID:   int64(u.IntID),
		Protected: protected,
		ContentID: key.IntID(),
	}

	newkey, err := retrievable.PlaceEntity(ctx, int64(0), &NewNote)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}
	log.Infof(ctx, "Information being submitted: ", NewNote, NewContent)
	http.Redirect(res, req, "/view/"+strconv.FormatInt(newkey.IntID(), 10), http.StatusSeeOther)
}

/// TODO: implement
func NOTES_GET_View(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	u, err := GetUserFromSession(req) // Check if a user is already logged in.
	ctx := appengine.NewContext(req)

	NoteKeyStr := params.ByName("ID")
	NoteKey, err := strconv.ParseInt(NoteKeyStr, 10, 64)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	ViewNote := &Note{}
	ViewContent := &Content{}

	err = retrievable.GetEntity(ctx, NoteKey, ViewNote)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	err = retrievable.GetEntity(ctx, ViewNote.ContentID, ViewContent)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	owner, err := GetUserFromID(ctx, ViewNote.OwnerID)
	if ErrorPage(ctx, res, nil, "Internal Server Error (x)", err, http.StatusSeeOther) {
		return
	}

	tempcontent := parse(ViewContent.Content)
	Body := template.HTML(tempcontent)

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
		Notekey:       NoteKeyStr,
		Content:       Body,
		User:          u,
		Owner:         owner,
	})

}

/// TODO: implement
func NOTES_GET_Editor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	u, err := GetUserFromSession(req) // Check if a user is already logged in.
	ctx := appengine.NewContext(req)

	NoteKeyStr := params.ByName("ID")
	NoteKey, err := strconv.ParseInt(NoteKeyStr, 10, 64)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	ViewNote := &Note{}
	ViewContent := &Content{}

	err = retrievable.GetEntity(ctx, NoteKey, ViewNote)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

	// Permission Check, For collaberation it can also check against a collaborator container after the user check.
	// When setting for example, privacy setting might only be able to be set by the Owner so a separation is still needed.
	if ViewNote.OwnerID != int64(u.IntID) && ViewNote.Protected {
		// Soft rejection. can also be substituted for a http Not Allowed.
		http.Redirect(res, req, "/view/"+NoteKeyStr, http.StatusSeeOther)
		return
	}

	err = retrievable.GetEntity(ctx, ViewNote.ContentID, ViewContent)
	if ErrorPage(ctx, res, nil, "Internal Server Error (2)", err, http.StatusSeeOther) {
		return
	}

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
		Notekey:       NoteKeyStr,
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

	tempcontent := parse(data)

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