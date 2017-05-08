package main

import (
	"bytes"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"math/rand"
	"mime/multipart"
	"strings"

	"github.com/Esseh/notorious-dev/AUTH"
	"github.com/Esseh/notorious-dev/CLOUD"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/COOKIE"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/FORUM"
	"github.com/Esseh/notorious-dev/NOTES"
	"github.com/Esseh/notorious-dev/NOTIFICATION"
	"github.com/Esseh/notorious-dev/PM"
	"github.com/Esseh/notorious-dev/USERS"
	"github.com/Esseh/retrievable"
	"github.com/disintegration/imaging"
	humanize "github.com/dustin/go-humanize" // russross markdown parser
	"github.com/pkg/errors"
	appcontext "golang.org/x/net/context"
)

func init() {
	// Tie functions into template here with ... "functionName":theFunction,
	funcMap := template.FuncMap{
		"getAvatarURL":     CORE.GetAvatarURL,
		"getUser":          GetUserFromID,
		"humanize":         humanize.Time,
		"humanizeSize":     humanize.Bytes,
		"monthfromtime":    CORE.MonthFromTime,
		"yearfromtime":     CORE.YearFromTime,
		"dayfromtime":      CORE.DayFromTime,
		"findsvg":          CORE.FindSVG,
		"findtemplate":     CORE.FindTemplate,
		"inc":              CORE.Inc,
		"addCtx":           CORE.AddCtx,
		"getDate":          CORE.GetDate,
		"toInt":            CORE.ToInt,
		"getMod":           GetMod,
		"canEditNote":      NOTES.CanEditNote,
		"canViewNote":      NOTES.CanViewNote,
		"findCollabs":      FindCollaborators,
		"canEdit":          NOTES.CanEditNote,
		"canView":          NOTES.CanViewNote,
		"getEmail":         GetEmail,
		"retrieveMessages": PM.RetrieveMessages,
		"getPageNumbers":   PM.GetPageNumbers,
		"incPage":          IncPage,
		"decPage":          DecPage,
		// "isOwner":       isOwner,
		"parse":            CORE.EscapeString,
		"getSubscriptions": NOTES.GetSubscriptions,
		"getNotifications": NOTIFICATION.GetNotifications,
		"isAdmin":          FORUM.IsAdmin,
		"getCategories":    FORUM.GetCategories,
		"getForums":        FORUM.GetForums,
		"getThreads":       FORUM.GetThreads,
		"getPosts":         FORUM.GetPosts,
	} // Load up all templates.
	CORE.TPL = template.New("").Funcs(funcMap)
	CORE.TPL = template.Must(CORE.TPL.ParseGlob("templates/*"))
}

type CropBounds struct {
	X         int
	Y         int
	W         int
	H         int
	RotateDeg int
}

func DecPage(i int64) int64 {
	if i == 0 {
		return 0
	}
	return i - 1
}
func IncPage(i int64) int64 { return i + 1 }

func GetMod(a int64) int64 {
	rand.Seed(a)
	return int64(rand.Uint32()) % 10
}

func GetEmail(ctx appcontext.Context, userID int64) string {
	u := USERS.User{}
	retrievable.GetEntity(ctx, userID, &u)
	return u.Email
}

func FindCollaborators(ctx CONTEXT.Context, c string) []int64 {
	dupCheck := make(map[int64]bool)
	temp := strings.Split(c, ":")
	var collabs []int64
	for _, x := range temp {
		nextEntry := AUTH.EmailReference{}
		err := retrievable.GetEntity(ctx, x, &nextEntry)
		if (err == nil) && (!dupCheck[nextEntry.UserID]) {
			collabs = append(collabs, nextEntry.UserID)
			dupCheck[nextEntry.UserID] = true
		}
	}
	return collabs
}

func GetUserFromID(ctx appcontext.Context, id int64) (*USERS.User, error) {
	owner := &USERS.User{}
	err := retrievable.GetEntity(ctx, id, owner)
	if err != nil {
		return &USERS.User{}, err
	} else {
		return owner, nil
	}
}

// Constructs the header.
// As the header gets more complex(such as capturing the current path)
// the need for such a helper function increases.
func MakeHeader(ctx CONTEXT.Context) *CONTEXT.HeaderData {
	oldCookie, err := COOKIE.GetValue(ctx.Req, "session")
	if err == nil {
		COOKIE.Make(ctx.Res, "session", oldCookie)
	}
	redirectURL := ctx.Req.URL.Path[1:]
	if redirectURL == "login" || redirectURL == "register" || redirectURL == "elevatedlogin" {
		redirectURL = ctx.Req.URL.Query().Get("redirect")
	}
	return &CONTEXT.HeaderData{
		ctx, ctx.User, redirectURL,
	}
}

//resize
func resizeImage(imgRdr io.Reader, x, y, width, height, finalWidth, finalHeight, rotateDeg int) (io.Reader, error) {
	img, _, err := image.Decode(imgRdr)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse image")
	}

	var rotated *image.NRGBA
	switch rotateDeg {
	case 90:
		rotated = imaging.Rotate270(img)
	case 180:
		rotated = imaging.Rotate180(img)
	case 270:
		rotated = imaging.Rotate90(img)
	default:
		rotated = imaging.Clone(img)
	}

	cropRect := image.Rect(x, y, x+width, y+height).Add(rotated.Bounds().Min)
	cropped := imaging.Crop(rotated, cropRect)

	resized := imaging.Resize(cropped, finalWidth, finalHeight, imaging.Lanczos)

	buf := &bytes.Buffer{}
	err = jpeg.Encode(buf, resized, nil)
	return buf, errors.Wrap(err, "Unable to convert to jpeg")
}

//upload
func uploadImage(ctx CONTEXT.Context, userID int64, header *multipart.FileHeader, cb *CropBounds, avatarReader io.ReadSeeker) error {
	imgRdr, err := resizeImage(avatarReader, cb.X, cb.Y, cb.W, cb.H, 500, 500, cb.RotateDeg)
	if err != nil {
		return errors.Wrap(err, "Resize image")
	}
	filename := CORE.GetAvatarPath(userID)
	return CLOUD.AddFile(ctx, filename, header.Header["Content-Type"][0], imgRdr)
}
