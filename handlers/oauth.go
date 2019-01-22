package handlers

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"context"
	"log"
	"encoding/base64"
	"crypto/rand"
	"encoding/json"
	"time"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/mchmarny/gauther/stores"
	"github.com/mchmarny/gauther/utils"

)

const (
	googleOAuthURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	stateCookieName = "authstate"
	userIDCookieName = "uid"
)

var (
	longTimeAgo    = time.Duration(3650 * 24 * time.Hour)
	cookieDuration = time.Duration(30 * 24 * time.Hour)
	oauthConfig *oauth2.Config
)


func getOAuthConfig(r *http.Request) *oauth2.Config {

	if oauthConfig != nil {
		return oauthConfig
	}

	// HTTPS or HTTP
	proto := r.Header.Get("X-FORWARDED-PROTO")
	if proto == "" {
		proto = "http"
	}

	baseURL := fmt.Sprintf("%s://%s", proto, r.Host)
	log.Printf("External URL: %s", baseURL)

	// OAuth
	oauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/callback", baseURL),
		ClientID:     utils.MustGetEnv("OAUTH_CLIENT_ID", ""),
		ClientSecret: utils.MustGetEnv("OAUTH_CLIENT_SECRET", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return oauthConfig

}



// OAuthLoginHandler handles oauth login
func OAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	uid := getCurrentUserID(r)
	if uid != "" {
		log.Printf("User ID from previous visit: %s", uid)
	}

	u := getOAuthConfig(r).AuthCodeURL(generateStateOauthCookie(w))
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// OAuthCallbackHandler handles oauth callback
func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {

	oauthState, _ := r.Cookie(stateCookieName)

	// checking state of the callback
	if r.FormValue("state") != oauthState.Value {
		err := errors.New("invalid oauth state from Google")
		ErrorHandler(w, r, err, http.StatusInternalServerError)
		return
	}

	// parsing callback data
	data, err := getOAuthedUserData(r)
	if err != nil {
		log.Printf("Error while parsing user data %v", err)
		ErrorHandler(w, r, err, http.StatusInternalServerError)
		return
	}

	dataMap := make(map[string]interface{})
	json.Unmarshal(data, &dataMap)

	email := dataMap["email"]
	log.Printf("Email: %s", email)
	id := utils.MakeID(email.(string))

	//server resize image
	pic := dataMap["picture"]
	if pic != nil {
		dataMap["picture"] = utils.ServerSizeResizePlusPic(pic.(string), 200)
	}


	// save data
	err = stores.SaveData(r.Context(), id, dataMap)
	if err != nil {
		log.Printf("Error while saving data: %v", err)
		ErrorHandler(w, r, err, http.StatusInternalServerError)
		return
	}

	// last auth date
	dataMap["last_auth"] = time.Now()


	// set cookie for 30 days
	cookie := http.Cookie{
		Name: userIDCookieName,
		Path: "/",
		Value: id,
		Expires: time.Now().Add(cookieDuration),
	}
	http.SetCookie(w, &cookie)

	// redirect on success
	http.Redirect(w, r, "/", http.StatusSeeOther)

}



// LogOutHandler resets cookie and redirects to home page
func LogOutHandler(w http.ResponseWriter, r *http.Request) {

	uid := getCurrentUserID(r)
	log.Printf("User logging out: %s", uid)

	cookie := http.Cookie{
		Name: userIDCookieName,
		Path: "/",
		Value: "",
		MaxAge: -1,
		Expires: time.Now().Add(-longTimeAgo),
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther) // home
}

func generateStateOauthCookie(w http.ResponseWriter) string {

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name: stateCookieName,
		Value: state,
		Expires: time.Now().Add(cookieDuration),
	}
	http.SetCookie(w, &cookie)

	return state
}

func getOAuthedUserData(r *http.Request) ([]byte, error) {

	// exchange code
	token, err := getOAuthConfig(r).Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return nil, fmt.Errorf("Got wrong exchange code: %v", err)
	}

	// user info
	response, err := http.Get(googleOAuthURL + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Error getting user info: %v", err)
	}
	defer response.Body.Close()

	// parse body
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	return contents, nil
}