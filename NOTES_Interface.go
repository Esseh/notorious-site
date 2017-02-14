package main
import (
	"errors"
	"strconv"
	"github.com/Esseh/retrievable"
	"google.golang.org/appengine/datastore"
)
// Retrieves an existing note and it's content by it's id.
func NOTES_GetExistingNote(ctx Context,id string)(*NOTES_Note,*NOTES_Content,error){
	RetrievedNote := &NOTES_Note{}
	RetrievedContent := &NOTES_Content{}
	NoteKey, err := strconv.ParseInt(id, 10, 64)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, NoteKey, RetrievedNote)
	if err != nil { return RetrievedNote,RetrievedContent,err }
	err = retrievable.GetEntity(ctx, RetrievedNote.ContentID, RetrievedContent)
	return RetrievedNote,RetrievedContent,err
}

// Verifies that the currently logged in user is allowed to interact with the Note.
func NOTES_VerifyNotePermission(ctx Context, note *NOTES_Note) bool {
	redirect := strconv.FormatInt(note.OwnerID, 10)
	if note.OwnerID != int64(ctx.user.IntID) && note.Protected {
		ctx.Redirect("/view/"+redirect)
		return false
	}
	return true
}

// Given a Content and Note it will construct instances of each, tie them together in the database and provide their keys.
func NOTES_CreateNewNote(ctx Context,NewContent NOTES_Content,NewNote NOTES_Note) (*datastore.Key,*datastore.Key,error) {
	contentKey, err := retrievable.PlaceEntity(ctx, int64(0), &NewContent)
	if err != nil { return contentKey,&datastore.Key{},err }
	NewNote.ContentID = contentKey.IntID()
	noteKey, err := retrievable.PlaceEntity(ctx, int64(0), &NewNote)
	return contentKey,noteKey,err
}

// Updates a note and its content based on the given id.
func NOTES_UpdateNoteContent(ctx Context,id string,UpdatedContent NOTES_Content, UpdatedNote NOTES_Note) error {
	Note := NOTES_Note{}
	noteID, err := strconv.ParseInt(id, 10, 64); if err != nil { return err }
	err = retrievable.GetEntity(ctx, noteID, &Note); if err != nil { return err }
	validated := NOTES_VerifyNotePermission(ctx, &Note); if !validated { return errors.New("Permission Error: Not Allowed") }
	if Note.OwnerID == int64(ctx.user.IntID) { Note.Protected = UpdatedNote.Protected }	
	_, err = retrievable.PlaceEntity(ctx, noteID, &Note); if err != nil { return err }
	_, err = retrievable.PlaceEntity(ctx, Note.ContentID, &UpdatedContent); return err
}

// Gets all the notes made by the AUTH_User
func NOTES_GetAllNotes(ctx Context, userID int64) ([]NOTES_NoteOutput, error) {
	q := datastore.NewQuery(NoteTable).Filter("OwnerID =", userID)
	res := []NOTES_Note{}
	output := make([]NOTES_NoteOutput, 0)
	keys, err := q.GetAll(ctx, &res)
	if err != nil {
		return nil, err
	}
	for i, _ := range res {
		var c NOTES_Content
		err := retrievable.GetEntity(ctx, res[i].ContentID, &c)
		if err != nil {
			return nil, err
		}
		output = append(output, NOTES_NoteOutput{keys[i].IntID(), res[i], c})
	}
	return output, nil
}