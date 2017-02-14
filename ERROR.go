package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/appengine/log"
)

var (
	ERROR_PasswordMatch           = errors.New("Passwords fields do not match")
	ERROR_InvalidLogin            = errors.New("Error Login: Parameters may not be nil!")
	ERROR_NotLoggedIn             = errors.New("Login Error: Cannot verify session!")
	ERROR_UsernameExists          = errors.New("Validation Error: Username already exists!")
	ERROR_InvalidPermission       = errors.New("Permission Error: Not Allowed")
	ERROR_NotImplemented          = errors.New("Structure Error: Function Not Implemented!")
	ERROR_EmptyField              = errors.New("Input Error: Required Field Empty")
	ERROR_NoUser                  = errors.New("Login Error: No Such User")
	ERROR_NoSession               = errors.New("Session does not exist")
	ERROR_TooLarge                = errors.New("Image Dimensions too large")
	ERROR_InvalidEmail            = errors.New("The email sent is invalid")
	ERROR_CategoryDoesNotExist    = errors.New("Category does not exist")
	ERROR_SubCategoryDoesNotExist = errors.New("Sub-Category does not exist")
	ERROR_MustOwnNotes            = errors.New("User does not own note")
	ERROR_MustOwnUpload           = errors.New("User does not own upload")
	ERROR_MustOwnItem             = errors.New("User does not own item")
	ERROR_NotMatchingHMac         = errors.New("Hmac checking failed")
	ERROR_NoItemToDelete          = errors.New("Must specify item to delete")
	ERROR_ItemDoesNotExist        = errors.New("Item does not exist")
)

/// Prints an error page to response and returns a boolean representation of the function executing.
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
/// Usage: Use if there is no constructive alternative.
func ERROR_Page(ctx Context, ErrorTitle string, e error, errCode int) bool {
	if e != nil {
		log.Errorf(ctx, "%s ---- %v\n", ErrorTitle, e)
		if ctx.user == nil {
			ctx.user = &User{}
		}
		args := &struct {
			Header    HeaderData
			ErrorName string
			ErrorDump error
			ErrorCode int
		}{
			HeaderData{ctx, ctx.user, ""}, ErrorTitle, e, errCode,
		}
		ctx.res.WriteHeader(errCode)
		ServeTemplateWithParams(ctx.res, "site-error", args)
		return true
	}
	return false
}

/// Returns to GET responding with FormValue("ErrorResponse")
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
/// Usage: Use in POST calls accessed from a GET of the same handle.
func ERROR_Back(ctx Context, err error, errorString string) bool {
	if err != nil {
		path := strings.Replace(ctx.req.URL.Path, "%2f", "/", -1)
		path += "?"+url.QueryEscape("ErrorResponse")+"="+url.QueryEscape(errorString)
		if ctx.req.FormValue("redirect") != "" {
			path += "&"+url.QueryEscape("redirect")+"="+ctx.req.FormValue("redirect")		
		}
		http.Redirect(ctx.res, ctx.req, path, http.StatusSeeOther)
		return true
	}
	return false
}

