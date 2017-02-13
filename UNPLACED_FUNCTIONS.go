package main
import (
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"strconv"
	"github.com/Esseh/retrievable"
)


type Context struct {
	context.Context
}

func NewContext(res http.ResponseWriter, req *http.Request) Context{
	return Context { appengine.NewContext(req) }
}

func GetNoteData(id string, ctx context.Context)(*Note,*Content,error){
	RetrievedNote := &Note{}
	RetrievedContent := &Content{}
	NoteKey, err := strconv.ParseInt(id, 10, 64)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, NoteKey, RetrievedNote)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, RetrievedNote.ContentID, RetrievedContent)
	return RetrievedNote,RetrievedContent,err
}

func VerifyNotePermission(res http.ResponseWriter, req *http.Request, user *User, note *Note) bool {
	if note.OwnerID != int64(user.IntID) && note.Protected {
		http.Redirect(res, req, "/view/"+params.ByName("ID"), http.StatusSeeOther)
		return false
	}
	return true
}