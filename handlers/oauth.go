package handlers

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"fmt"
	"io/ioutil"
	"context"
	"strings"
	"log"
	"encoding/base64"
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/mchmarny/gauther/stores"
	"github.com/mchmarny/gauther/utils"
)

const (
	googleOAuthURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	stateCookieName = "authstate"
	userIDCookieName = "uid"
)

var (
	oauthConfig *oauth2.Config
)

// ConfigureOAuthHandler initializes auth handler
func ConfigureOAuthHandler(ctx context.Context, baseURL string) error {
	if baseURL == "" || !strings.HasPrefix(baseURL, "http") {
		return fmt.Errorf("baseURL must start with HTTP or HTTPS")
	}

	log.Printf("Configuring auth callback to %s", baseURL)
	oauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/callback", baseURL),
		ClientID:     utils.MustGetEnv("OAUTH_CLIENT_ID", ""),
		ClientSecret: utils.MustGetEnv("OAUTH_CLIENT_SECRET", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return stores.InitStore(ctx)
}

// OAuthLoginHandler handles oauth login
func OAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	uidCookie, _ := r.Cookie(userIDCookieName)
	if uidCookie != nil {
		log.Printf("User authenticated: %s", uidCookie.Value)
	}
	u := oauthConfig.AuthCodeURL(generateStateOauthCookie(w))
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// OAuthCallbackHandler handles oauth callback
func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {

	oauthState, _ := r.Cookie(stateCookieName)

	// checking state of the callback
	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth state from Google")
		ErrorHandler(w, r)
		return
	}

	// parsing callback data
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Printf("Error while parsing user data %v", err)
		ErrorHandler(w, r)
		return
	}

	dataMap := make(map[string]interface{})
	json.Unmarshal(data, &dataMap)

	email := dataMap["email"]
	log.Printf("Email: %s", email)
	id := utils.MakeID(email.(string))


	// save data
	err = stores.SaveData(r.Context(), id, dataMap)
	if err != nil {
		log.Printf("Error while saving data: %v", err)
		ErrorHandler(w, r)
		return
	}

	// set cookie for 30 days
	cookie := http.Cookie{
		Name: userIDCookieName,
		Path: "/",
		Value: id,
		Expires: time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	// redirect on success
	http.Redirect(w, r, "/", http.StatusSeeOther)

}



// LogOutHandler resets cookie and redirects to home page
func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name: userIDCookieName,
		Path: "/",
		Value: "",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func generateStateOauthCookie(w http.ResponseWriter) string {
	exp := time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name: stateCookieName,
		Value: state,
		Expires: exp,
	}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {

	// exchange code
	token, err := oauthConfig.Exchange(context.Background(), code)
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