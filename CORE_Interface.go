package main
import (
	"strings"
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)
// Constructs an instance of Context
func NewContext(res http.ResponseWriter, req *http.Request) Context{
	user, err := AUTH_GetUserFromSession(req)
	ctx := Context { 
		req: req,
		res: res,
		user: user,
		userException: err,
	}
	ctx.Context = appengine.NewContext(req)
	return ctx
}

// A black box that automatially keeps track of transaction timing for the database
// and stores useful metadata.
type Context struct {
	// The active request 
	req *http.Request
	// The output writer to the user's browser.
	res http.ResponseWriter
	// The currently logged in user.
	user *USER_User
	// Any problems that occured while logging in.
	userException error
	// transaction timing information.
	context.Context
}

// Returns true if the user is not logged in.
func (ctx Context)AssertLoggedInFailed() bool {
	if ctx.userException != nil {
		path := strings.Replace(ctx.req.URL.Path[1:], "%2f", "/", -1)
		http.Redirect(ctx.res, ctx.req, PATH_AUTH_Login+"?redirect="+path, http.StatusSeeOther)
		return true
	}
	return false
}

// Simplified redirect, useful for general redirects. If the redirect demands a more severe status code use tradition http.Redirect.
func (ctx Context)Redirect(uri string){ http.Redirect(ctx.res, ctx.req, uri, http.StatusSeeOther) }