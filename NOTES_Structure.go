package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	NoteTable = "NotePermissions"
)

// Contains the metadata of a Note and represents the existance of one.
type NOTES_Note struct {
	// The ID to the AUTH_User owner.
	OwnerID int64
	// A boolean value representing whether or not it can be publically edited.
	Protected bool
	// The ID to the related NOTES_Content
	ContentID int64
}

// The actual content referred to by NOTES_Note
type NOTES_Content struct {
	// The title and content respectively.
	Title, Content string
}

func (n *NOTES_Note) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, NoteTable, "", key.(int64), nil)
}
func (n *NOTES_Content) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "NoteContents", "", key.(int64), nil)
}


// A convience structure used with GetAllNotes, provides a clean managable output on the front end.
type NOTES_NoteOutput struct {
	// The ID for a note.
	ID      int64
	// The Note Meta Data
	Data    NOTES_Note
	// The Note Content
	Content NOTES_Content
}