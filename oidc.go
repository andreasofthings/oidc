package main

// [START import]
import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var (
	ctx              context.Context
	conf             *oauth2.Config
	oauthStateString = "pseudo-random"
)

// https://blog.kowalczyk.info/article/f/accessing-github-api-from-go.html

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var htmlIndex = `<html>
	<body>
		<a href="/login">Log In</a>
	</body>
	</html>`
	w.Write([]byte(htmlIndex))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	// loginHandler initiates an OAuth flow to authenticate the user.
	url := conf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// /callback Verify oauth2 state and errors.
	var state = r.FormValue("state")
	var code = r.FormValue("code")

	if state != oauthStateString {
		w.Write([]byte("Invalid State"))
		return // "invalid oauth state"
	}

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		w.Write([]byte("Code Exchange error"))
		return // "code exchange failed" //  %s", err.Error())
	}

	client := conf.Client(ctx, token)

	response, err := client.Get("https://www.pramari.de/api/user")
	if err != nil {
		w.Write([]byte("Generic API error"))
		return // "failed getting user info" // : %s", err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.Write([]byte("Read API error"))
		return // "failed getting user info" // : %s", err.Error())
	}
	w.Write(body)
}

func main() {
	ctx = context.Background()

	conf = &oauth2.Config{
		ClientID:     "Y5Uy9Cg6vpaYE0bOyjraS52JoNw8Z3BKiCdLl1k1",
		ClientSecret: "Sw4n6OP54IY18e0wBOMc4vKcZgH2ADJxKSSl0vV9Pg04znLbGJYaMPtPoL6JqFZK2UVlMj30toaJEhSW76xn7NTqRZjC1NFw3bE2vO0iMnsH2fr2xAicpr0XsQYABJc9",
		Scopes:       []string{"read", "write"},
		RedirectURL:  "http://localhost:8000/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.pramari.de/oauth2/authorize",
			TokenURL: "https://www.pramari.de/oauth2/token",
		},
	}

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/callback", callbackHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}
