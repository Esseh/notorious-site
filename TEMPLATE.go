package main

import (
	"html/template"
	humanize "github.com/dustin/go-humanize" // russross markdown parser
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/COOKIE"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/USERS"
)



func init() {
	// Tie functions into template here with ... "functionName":theFunction,
	funcMap := template.FuncMap{
		"getAvatarURL":  CORE.GetAvatarURL,
		"getUser":       USERS.GetUserFromID,
		"humanize":      humanize.Time,
		"humanizeSize":  humanize.Bytes,
		"monthfromtime": CORE.MonthFromTime,
		"yearfromtime":  CORE.YearFromTime,
		"dayfromtime":   CORE.DayFromTime,
		"findsvg":       CORE.FindSVG,
		"findtemplate":  CORE.FindTemplate,
		"inc":           CORE.Inc,
		"addCtx":        CORE.AddCtx,
		"getDate":       CORE.GetDate,
		"toInt":		 CORE.ToInt,
		// "isOwner":       isOwner,
		"parse": CORE.EscapeString,
	} // Load up all templates.
	CORE.TPL = template.New("").Funcs(funcMap)
	CORE.TPL = template.Must(CORE.TPL.ParseGlob("templates/*"))

}

// Constructs the header.
// As the header gets more complex(such as capturing the current path)
// the need for such a helper function increases.
func MakeHeader(ctx CONTEXT.Context) *CONTEXT.HeaderData {
	oldCookie, err := COOKIE.GetValue(ctx.Req, "session")
	if err == nil { COOKIE.Make(ctx.Res, "session", oldCookie) }
	redirectURL := ctx.Req.URL.Path[1:]
	if redirectURL == "login" || redirectURL == "register" || redirectURL == "elevatedlogin" {
		redirectURL = ctx.Req.URL.Query().Get("redirect")
	}
	return &CONTEXT.HeaderData{
		ctx, ctx.User, redirectURL,
	}
}

