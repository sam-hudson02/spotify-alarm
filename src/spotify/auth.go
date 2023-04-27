package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func readConfig() {
	godotenv.Load()
}

func getAuth() *spotifyauth.Authenticator {
	fmt.Println("Reading config...")
	readConfig()
	auth := spotifyauth.New(spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT")), spotifyauth.WithScopes(spotifyauth.ScopeUserModifyPlaybackState, spotifyauth.ScopeUserLibraryRead))
	return auth
}

func GetUrl() string {
	auth := getAuth()
	url := auth.AuthURL("123")
	return url
}

func cacheToken(token *oauth2.Token) {
	fmt.Println("Caching token...")
	b, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(".cache", b, 0644)
}

func ParseToken() (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(".cache")
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(b, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func UseToken(token *oauth2.Token) *spotify.Client {
	auth := getAuth()
	webClient := auth.Client(context.Background(), token)
	client := spotify.New(webClient)
	return client
}

func ParseCode(r *http.Request) (*spotify.Client, error) {
	auth := getAuth()
	tok, err := auth.Token(r.Context(), "123", r)
	cacheToken(tok)

	if err != nil {
		return nil, err
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	return client, nil
}
