// Contains the handlers for our OAuth interactions.
package main
import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/Esseh/notorious-dev/AUTH"
	"github.com/Esseh/notorious-dev/PATHS"
	"github.com/Esseh/goauth"
	github "github.com/Esseh/goauth-github"
	dropbox "github.com/Esseh/goauth-dropbox"
)

func init(){
	goauth.GlobalSettings.ClientType = "appengine"
	github.Config.ClientID = "e0297346f88565c9f443"
	github.Config.SecretID = "7dd96d4a262a004aeffefe4b0af1a38e03b38d14"
	github.Config.Redirect = "http://localhost:8080/login/github/oauth/recieve"
	dropbox.Config.ClientID = "ddhu8e7nswl56yt"
	dropbox.Config.SecretID = "387kru0n9nb0qkk"
	dropbox.Config.Redirect = "http://localhost:8080/login/dropbox/oauth/recieve"
}

func INIT_OAUTH_Handlers(r *httprouter.Router){
	r.GET(PATHS.AUTH_OAUTH_GITHUB_Send, AUTH_OAUTH_GITHUB_Send)
	r.GET(PATHS.AUTH_OAUTH_GITHUB_Recieve, AUTH_OAUTH_GITHUB_Recieve)
	r.GET(PATHS.AUTH_OAUTH_DROPBOX_Send, AUTH_OAUTH_DROPBOX_Send)
	r.GET(PATHS.AUTH_OAUTH_DROPBOX_Recieve, AUTH_OAUTH_DROPBOX_Recieve)
}
	
func AUTH_OAUTH_GITHUB_Send(res http.ResponseWriter, req *http.Request, params httprouter.Params) { 
	github.Send(res,req)
}

func AUTH_OAUTH_GITHUB_Recieve(res http.ResponseWriter, req *http.Request, params httprouter.Params){
	token := github.Recieve(res,req)
	emailObj, _ := token.Email(req)
	accountObj, _ := token.AccountInfo(req)
	AUTH.OAuthLogin(req, res, emailObj.Email, accountObj.Login, "", token.State)
}

func AUTH_OAUTH_DROPBOX_Send(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	dropbox.Send(res,req)
}
func AUTH_OAUTH_DROPBOX_Recieve(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	token := dropbox.Recieve(res,req)
	accountInfo, _ := token.AccountInfo(req)
	AUTH.OAuthLogin(req, res, accountInfo.Email, accountInfo.NameDetails.GivenName, accountInfo.NameDetails.Surname, token.State)	
}