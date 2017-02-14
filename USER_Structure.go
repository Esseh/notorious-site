// Contains the object structure and methods relating to the USER.
package main

import (
	"encoding/json"

	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	UsersTable          = "Users"
	RecentlyViewedTable = "RecentlyViewedCourses"
)

type (
	// Represents an individual user.
	USER_User struct {
		// First and Last Name
		First, Last       string
		Email             string
		// Whether they have an active avatar or not.
		Avatar            bool `datastore:",noindex"`
		// Biography.
		Bio               string
		// ID referred to itself.
		retrievable.IntID `datastore:"-" json:"-"`
	}
	// An encrypted user.
	USER_EncryptedUser struct {
		First, Last string
		Email       string
		Avatar      bool `datastore:",noindex"`
		Bio         string
	}
)

func (u *USER_User) Key(ctx context.Context, key interface{}) *datastore.Key {
	if v, ok := key.(retrievable.IntID); ok {
		return datastore.NewKey(ctx, UsersTable, "", int64(v), nil)
	}
	return datastore.NewKey(ctx, UsersTable, "", key.(int64), nil)
}

// Converts user to an encrypted user.
func (u *USER_User) toEncrypt() (*USER_EncryptedUser, error) {
	e := &USER_EncryptedUser{
		First:     u.First,
		Last:      u.Last,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
	}
	email, err := AUTH_Encrypt([]byte(u.Email), encryptKey)
	if err != nil { return nil, err }
	e.Email = email
	return e, nil
}

// Converts encrypted user to normal user.
func (u *USER_User) fromEncrypt(e *USER_EncryptedUser) error {
	email, err := AUTH_Decrypt(e.Email, encryptKey)
	if err != nil { return err }
	u.First = e.First
	u.Last = e.Last
	u.Email = string(email)
	u.Avatar = e.Avatar
	u.Bio = e.Bio
	return nil
}

// User -> JSON
func (u *USER_User) Serialize() []byte {
	data, _ := u.toEncrypt()
	ret, _ := json.Marshal(&data)
	return ret
}
// JSON -> User
func (u *USER_User) Unserialize(data []byte) error {
	e := &USER_EncryptedUser{}
	err := json.Unmarshal(data, e)
	if err != nil { return err }
	return u.fromEncrypt(e)
}

