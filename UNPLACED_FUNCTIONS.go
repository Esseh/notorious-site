package main
import (
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)


type Context struct {
	context.Context
}

func NewContext(res http.ResponseWriter, req *http.Request) Context{
	return Context { appengine.NewContext(req) }
}