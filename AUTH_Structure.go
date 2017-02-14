/// Contains the structure for our logins/user structure.
package main

import (
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

const (
	LoginTable         = "Login"
	SessionTable       = "Session"
	PasswordResetTable = "PasswordReset"
)

type (
	// Contains local authentication information about the user.
	AUTH_LoginLocalAccount struct {
		// Key for USER_User
		UserID   int64
		// An encrypted password.
		Password []byte
	}
	// Contains session information from an individual login instance.
	AUTH_Session struct {
		// Key for USER_User
		UserID       int64
		// Key for local instance of this object
		ID           int64  `datastore:"-"`
		// IP Address at point of login.
		IP           string `datastore:",noindex"`
		// Browser Information at point of login.
		BrowserUsed  string `datastore:",noindex"`
		// Physical Location at point of login.
		LocationUsed string `datastore:",noindex"`
		// Time that the session was created.
		LastUsed     time.Time
	}
)

// Retrieves the Local Login Account, overloaded to handle OAUTH.
func (l *AUTH_LoginLocalAccount) Get(ctx context.Context, key interface{}) error {
	getErr := retrievable.GetEntity(ctx, key, l) 
	if getErr != nil {                           
		oauth := LoginOauthAccount{}
		ogetErr := retrievable.GetEntity(ctx, key, &oauth)
		if ogetErr != nil { return ogetErr }
		l.UserID = oauth.UserID
	}
	return nil
}

//Updates the Local Login Account, overloaded to handle OAUTH.
func (l *AUTH_LoginLocalAccount) Place(ctx context.Context, key interface{}) (*datastore.Key, error) {
	if string(l.Password) == "" { // OAuth Case
		oauth := LoginOauthAccount{}
		oauth.UserID = oauth.UserID
		return retrievable.PlaceEntity(ctx, key, &oauth)
	} else { // LoginLocal Case
		return retrievable.PlaceEntity(ctx, key, l)
	}
}

func (l *AUTH_LoginLocalAccount) Key(ctx context.Context, key interface{}) *datastore.Key {
	e, _ := AUTH_Encrypt([]byte(key.(string)), encryptKey)
	return datastore.NewKey(ctx, LoginTable, e, 0, nil)
}
func (s *AUTH_Session) Key(ctx context.Context, key interface{}) *datastore.Key {
	return datastore.NewKey(ctx, SessionTable, "", key.(int64), nil)
}

func (s *AUTH_Session) StoreKey(key *datastore.Key) { s.ID = key.IntID() }
