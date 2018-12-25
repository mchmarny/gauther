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
	"os"
	"time"

	"github.com/mchmarny/gauther/stores"
)

const googleOAuthURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var oauthConfig *oauth2.Config

// ConfigureOAuthHandler initializes auth
func ConfigureOAuthHandler(ctx context.Context, baseURL string) error {

	if baseURL == "" || !strings.HasPrefix(baseURL, "http") {
		return fmt.Errorf("baseURL must start with HTTP or HTTPS")
	}

	log.Printf("Configuring auth callback to %s", baseURL)

	oauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/callback", baseURL),
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" {
		log.Fatalf("Both OAUTH_CLIENT_ID and OAUTH_CLIENT_SECRET must be defined.")
	}

	return stores.InitStore(ctx)

}

// OAuthLoginHandler handles oauth login
func OAuthLoginHandler(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
	AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
	You must always provide a non-empty string and validate that it matches the the
	state query parameter on your redirect callback.
	*/
	u := oauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// OAuthCallbackHandler handles oauth callback
func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("Raw data: %s", data)
	dataMap := make(map[string]interface{})
	json.Unmarshal(data, &dataMap)

	log.Printf("Data map: %v", dataMap)

	email := dataMap["email"]
	log.Printf("Email: %s", email)

	id := stores.MakeID(email.(string))
	log.Printf("ID: %s", id)

	err = stores.SaveData(r.Context(), id, dataMap)
	if err != nil {
		log.Printf("Error while saving data: %v", err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// TODO: redirect to landing page
	fmt.Fprintf(w, "%s", data)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {

	// Use code to get token and get user info from Google.
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(googleOAuthURL + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	return contents, nil
}