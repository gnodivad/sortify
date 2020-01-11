package api

import (
	"gnodivad/sortify/src/config"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func StartAuth(w http.ResponseWriter, r *http.Request) {
	conf := config.Init()
	auth.SetAuthInfo(conf.Spotify.ClientID, conf.Spotify.SecretKey)
	url := auth.AuthURL(state)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	conf := config.Init()
	auth.SetAuthInfo(conf.Spotify.ClientID, conf.Spotify.SecretKey)

	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
}
