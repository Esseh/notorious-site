package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	NoteTable = "NotePermissions"
)

type Note struct {
	OwnerID int64
	Protected bool
	ContentID int64
}

type Content struct {
	Title, Content string
}

func (n *Note) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, NoteTable, "", key.(int64), nil)
}
func (n *Content) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "NoteContents", "", key.(int64), nil)
}
