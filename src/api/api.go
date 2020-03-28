package api

import (
	"fmt"
	"gnodivad/sortify/src/config"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	client := auth.NewClient(tok)

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.ID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(conf.Jwt.SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "user-token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	username := getUser(r)
	fmt.Println(username)
}

func getUser(r *http.Request) string {
	conf := config.Init()
	c, err := r.Cookie("user-token")
	if err != nil {
		if err == http.ErrNoCookie {
			return ""
		}
		return ""
	}

	tknStr := c.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Jwt.SecretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Fatal(err)
			return ""
		}
		log.Fatal(err)
		return ""
	}
	if !token.Valid {
		log.Fatal("Token is invalid")
		return ""
	}

	return claims.Username
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
