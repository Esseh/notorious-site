package main
import (
	"strings"
	"errors"
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"strconv"
	"github.com/Esseh/retrievable"
	"google.golang.org/appengine/datastore"
)


type Context struct {
	req *http.Request
	res http.ResponseWriter
	user *User
	userException error
	context.Context
}

func (ctx Context)AssertLoggedInFailed() bool {
	if ctx.userException != nil {
		path := strings.Replace(ctx.req.URL.Path[1:], "%2f", "/", -1)
		http.Redirect(ctx.res, ctx.req, PATH_AUTH_Login+"?redirect="+path, http.StatusSeeOther)
		return true
	}
	return false
}

func (ctx Context)Redirect(uri string){ http.Redirect(ctx.res, ctx.req, uri, http.StatusSeeOther) }

func NewContext(res http.ResponseWriter, req *http.Request) Context{
	user, err := GetUserFromSession(req)
	ctx := Context { 
		req: req,
		res: res,
		user: user,
		userException: err,
	}
	ctx.Context = appengine.NewContext(req)
	return ctx
}

func GetExistingNote(ctx Context,id string)(*Note,*Content,error){
	RetrievedNote := &Note{}
	RetrievedContent := &Content{}
	NoteKey, err := strconv.ParseInt(id, 10, 64)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, NoteKey, RetrievedNote)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, RetrievedNote.ContentID, RetrievedContent)
	return RetrievedNote,RetrievedContent,err
}

func VerifyNotePermission(ctx Context, note *Note) bool {
	redirect := strconv.FormatInt(note.OwnerID, 10)
	if note.OwnerID != int64(ctx.user.IntID) && note.Protected {
		ctx.Redirect("/view/"+redirect)
		return false
	}
	return true
}

func CreateNewNote(ctx Context,NewContent Content,NewNote Note) (*datastore.Key,*datastore.Key,error) {
	contentKey, err := retrievable.PlaceEntity(ctx, int64(0), &NewContent)
	if err != nil { return contentKey,&datastore.Key{},err }
	NewNote.ContentID = contentKey.IntID()
	noteKey, err := retrievable.PlaceEntity(ctx, int64(0), &NewNote)
	return contentKey,noteKey,err
}

func UpdateNoteContent(ctx Context,id string,UpdatedContent Content, UpdatedNote Note) error {
	Note := Note{}
	noteID, err := strconv.ParseInt(id, 10, 64); if err != nil { return err }
	err = retrievable.GetEntity(ctx, noteID, &Note); if err != nil { return err }
	validated := VerifyNotePermission(ctx, &Note); if !validated { return errors.New("Permission Error: Not Allowed") }
	if Note.OwnerID == int64(ctx.user.IntID) { Note.Protected = UpdatedNote.Protected }	
	_, err = retrievable.PlaceEntity(ctx, noteID, &Note); if err != nil { return err }
	_, err = retrievable.PlaceEntity(ctx, Note.ContentID, &UpdatedContent); return err
}


func ValidLogin(username,password string) bool {
	return password != "" && username != ""
}

func LoginToWebsite(ctx Context,username,password string) (string, error) {
	userID, err := GetUserIDFromLogin(ctx, strings.ToLower(username), password)
	if err != nil { return "Login Information Is Incorrect", err }
	sessionID, err := CreateSessionID(ctx, ctx.req, userID)
	if err != nil { return "Login error, try again later.", err }
	err = MakeCookie(ctx.res, "session", strconv.FormatInt(sessionID, 10))
	return "Login error, try again later.",err
}

func LogoutFromWebsite(ctx Context)(string, error){
	sessionIDStr, err := GetCookieValue(ctx.req, "session")
	if err != nil { return "Must be logged in", err }
	sessionVal, err := strconv.ParseInt(sessionIDStr, 10, 0)	
	if err != nil { return "Bad cookie value", err }
	err = retrievable.DeleteEntity(ctx, (&Session{}).Key(ctx, sessionVal))
	if err == nil { DeleteCookie(ctx.res, "session") }
	return "No such session found!", err
}

func RegisterNewUser(ctx Context, username, password, confirmPassword, firstName, lastName string)(string,error){
	newUser := &User{ // Make the New User
		Email:    strings.ToLower(username),
		First:    firstName,
		Last:     lastName,
	}		
	if !ValidLogin(username,password) { return "Invalid Login Information", errors.New("Bad Login") }
	if password != confirmPassword { return "Passwords Do Not Match", errors.New("Password Mismatch") }
	_, err := CreateUserFromLogin(ctx, newUser.Email, password, newUser)
	return "Username Taken", err
}