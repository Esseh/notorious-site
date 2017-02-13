// Contains helper functions for dealing with authentication
package AUTH

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Esseh/retrievable"
	"github.com/mssola/user_agent"
	"github.com/pariz/gountries"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)



func AUTH_GetUserIDFromLogin(ctx context.Context, email, password string) (int64, error) {
	urID := AUTH_LoginLocalAccount{}
	if getErr := retrievable.GetEntity(ctx, email, &urID); getErr != nil { return -1, getErr }
	if compareErr := bcrypt.CompareHashAndPassword(urID.Password, []byte(password)); compareErr != nil {
		return -1, compareErr
	}
	return urID.UserID, nil
}

func AUTH_CreateUserFromLogin(ctx context.Context, email, password string, u *User) (*User, error) {
	checkLogin := AUTH_LoginLocalAccount{}
	// Check that user does not exist
	if checkErr := retrievable.GetEntity(ctx, email, &checkLogin); checkErr == nil {
		return u, ErrUsernameExists
	} else if checkErr != datastore.ErrNoSuchEntity && checkErr != nil {
		return u, checkErr
	}

	ukey, putUserErr := retrievable.PlaceEntity(ctx, retrievable.IntID(0), u)
	if putUserErr != nil { return u, putUserErr }
	if u.IntID == 0 { return u, errors.New("HEY, DATASTORE IS STUPID") }

	cryptPass, cryptErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if cryptErr != nil { return u, cryptErr }

	uLogin := AUTH_LoginLocalAccount{
		Password: cryptPass,
		UserID:   ukey.IntID(),
	}
	_, putErr := retrievable.PlaceEntity(ctx, email, &uLogin)
	return u, putErr
}

func AUTH_CreateSessionID(ctx context.Context, req *http.Request, userID int64) (sessionID int64, _ error) {
	agent := user_agent.New(req.Header.Get("user-agent"))
	browse, vers := agent.Browser()
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil { ip = req.RemoteAddr }
	country := req.Header.Get("X-AppEngine-Country")
	region := req.Header.Get("X-AppEngine-Region")
	city := req.Header.Get("X-AppEngine-City")
	location, err := AUTH_getLocationName(country, strings.ToUpper(region))
	if err != nil {
		location = "Unknown"
	} else {
		location = strings.Title(city) + ", " + location
	}
	newSession := AUTH_Session{
		UserID:      userID,
		BrowserUsed: browse + " " + vers,
		IP:          ip,
		// Test While Deployed
		LocationUsed: location,
		LastUsed:     time.Now(),
	}
	rk, err := retrievable.PlaceEntity(ctx, int64(0), &newSession)
	if err != nil { return int64(-1), err }
	return rk.IntID(), err
}

func AUTH_getLocationName(country, region string) (string, error) {
	c, err := gountries.New().FindCountryByAlpha(country)
	if err != nil { return "", err }
	for _, r := range c.SubDivisions() {
		if r.Code == region {
			return r.Name + ", " + c.Name.BaseLang.Common, nil
		}
	}
	return c.Name.BaseLang.Common, nil
}

func AUTH_GetUserIDFromSession(ctx context.Context, sessionID int64) (userID int64, _ error) {
	sessionData, err := AUTH_GetSession(ctx, sessionID)
	if err != nil { return 0, err }
	return sessionData.UserID, nil
}

func AUTH_GetSession(ctx context.Context, sessionID int64) (AUTH_Session, error) {
	s := AUTH_Session{}
	getErr := retrievable.GetEntity(ctx, sessionID, &s) // Get actual session from datastore
	if getErr != nil { return AUTH_Session{}, ErrNotLoggedIn }
	s.LastUsed = time.Now()
	if _, err := retrievable.PlaceEntity(ctx, sessionID, &s); err != nil { return AUTH_Session{}, err }
	return s, nil
}

func AUTH_GetSessionID(req *http.Request) (int64, error) {
	sessionIDStr, err := GetCookieValue(req, "session")
	if err != nil { return -1, ErrNotLoggedIn }
	id, err := strconv.ParseInt(sessionIDStr, 10, 64) // Change cookie val into key
	if err != nil { return -1, ErrInvalidLogin }
	return id, nil
}

func AUTH_GetUserFromSession(req *http.Request) (*User, error) {
	userID, err := AUTH_GetUserIDFromRequest(req)
	if err != nil { return &User{}, err }
	ctx := appengine.NewContext(req)
	return AUTH_GetUserFromID(ctx, userID)
}

func AUTH_GetUserIDFromRequest(req *http.Request) (int64, error) {
	s, err := AUTH_GetSessionID(req)
	if err != nil { return 0, err }
	ctx := appengine.NewContext(req)
	userID, err := AUTH_GetUserIDFromSession(ctx, s)
	if err != nil { return 0, err }
	return userID, nil
}

func AUTH_GetUserFromID(ctx context.Context, userID int64) (*User, error) {
	u := &User{}
	getErr := retrievable.GetEntity(ctx, retrievable.IntID(userID), u)
	return u, getErr
}

func AUTH_ValidLogin(username,password string) bool {
	return password != "" && username != ""
}

func AUTH_LoginToWebsite(ctx Context,username,password string) (string, error) {
	userID, err := AUTH_GetUserIDFromLogin(ctx, strings.ToLower(username), password)
	if err != nil { return "Login Information Is Incorrect", err }
	sessionID, err := AUTH_CreateSessionID(ctx, ctx.req, userID)
	if err != nil { return "Login error, try again later.", err }
	err = MakeCookie(ctx.res, "session", strconv.FormatInt(sessionID, 10))
	return "Login error, try again later.",err
}

func AUTH_LogoutFromWebsite(ctx Context)(string, error){
	sessionIDStr, err := GetCookieValue(ctx.req, "session")
	if err != nil { return "Must be logged in", err }
	sessionVal, err := strconv.ParseInt(sessionIDStr, 10, 0)	
	if err != nil { return "Bad cookie value", err }
	err = retrievable.DeleteEntity(ctx, (&AUTH_Session{}).Key(ctx, sessionVal))
	if err == nil { DeleteCookie(ctx.res, "session") }
	return "No such session found!", err
}

func AUTH_RegisterNewUser(ctx Context, username, password, confirmPassword, firstName, lastName string)(string,error){
	newUser := &User{ // Make the New User
		Email:    strings.ToLower(username),
		First:    firstName,
		Last:     lastName,
	}		
	if !AUTH_ValidLogin(username,password) { return "Invalid Login Information", errors.New("Bad Login") }
	if password != confirmPassword { return "Passwords Do Not Match", errors.New("Password Mismatch") }
	_, err := AUTH_CreateUserFromLogin(ctx, newUser.Email, password, newUser)
	return "Username Taken", err
}