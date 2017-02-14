// A set of functions for making and maintaining cookies.
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)
const( 
	sessionTime = 7 * 24 * 60 * 60
	HMAC_Key = "csci150project2016"
)
func createHmac(value string) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(HMAC_Key))
	_, err := io.WriteString(mac, value)
	if err != nil {
		return []byte{}, err
	}
	return mac.Sum(nil), nil
}

func splitMac(value string) (string, string) {
	i := strings.LastIndex(value, ".")
	if i == -1 {
		return value, ""
	}
	return value[:i], value[i+1:]
}

func checkMac(value, mac string) bool {
	derivedMac, err := createHmac(value)
	if err != nil {
		return false
	}
	macData, err := base64.RawURLEncoding.DecodeString(mac)
	if err != nil {
		return false
	}
	return hmac.Equal(derivedMac, macData)
}

// Deletes a cookie held in the current session by name.
func COOKIE_Delete(res http.ResponseWriter, name string) {
	http.SetCookie(res, &http.Cookie{
		Name:   name,
		MaxAge: -1,
		Path:   "/",
	})
}

// Initializes a cookie into the current session.
func COOKIE_Make(res http.ResponseWriter, name, value string) error {
	mac, err := createHmac(value)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     name,
		Value:    value + "." + base64.RawURLEncoding.EncodeToString(mac),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   sessionTime,
	}
	http.SetCookie(res, c)
	return nil
}

// Retrieves the value located inside of a cookie.
func COOKIE_GetValue(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}
	val, mac := splitMac(cookie.Value)
	if good := checkMac(val, mac); !good {
		return "", ERROR_NotMatchingHMac
	}
	return val, nil
}
