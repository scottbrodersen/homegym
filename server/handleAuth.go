package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/scottbrodersen/homegym/auth"
	"github.com/scottbrodersen/homegym/workoutlog"
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "use GET to log in", http.StatusMethodNotAllowed)
		return
	}

	credentials := new(creds)
	// GET means we received a form submission
	if r.Method == http.MethodGet {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "could not parse credentials", http.StatusBadRequest)
			return
		}

		credentials.Username = r.FormValue("username")
		credentials.Password = r.FormValue("password")

		if credentials.Username == "" || credentials.Password == "" {
			http.Error(w, "You must provide both a user name and a password", http.StatusBadRequest)
			return
		}
	} else {
		// POST means we received a fetch request
		*credentials = readPostedBody(w, r, 10000)
	}

	token, sessionID, err := authorizer.IssueToken(credentials.Username, credentials.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}
		log.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	tokenCookie := tokenCookieMaker(token)
	sessionCookie := &http.Cookie{
		Name:     "session",
		Value:    *sessionID,
		Secure:   isSafe,
		HttpOnly: isSafe,
		Path:     "/homegym/",
		SameSite: samesite,
		MaxAge:   authorizer.SessionTTL(),
	}

	http.SetCookie(w, tokenCookie)
	http.SetCookie(w, sessionCookie)
	h := w.Header()
	standardHeaders(&h)

	if r.Method == http.MethodGet {
		// redirect to home handler
		http.Redirect(w, r, "/homegym/home/", http.StatusFound)

	}
	w.WriteHeader(http.StatusOK)
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "request must be a POST", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Error(err)
		http.Error(w, "could not parse user info", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || password == "" || email == "" {
		http.Error(w, "missing required field values", http.StatusBadRequest)
		return
	}
	// TODO: check if this is the first user and if so make them admin
	role := auth.User

	_, err := workoutlog.FrontDesk.NewUser(username, role, email, password)

	if err != nil {
		log.Error(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     cookieUsername,
		Value:    username,
		Secure:   isSafe,
		HttpOnly: isSafe,
		Path:     "/homegym/login/",
		SameSite: samesite,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/homegym/login/", http.StatusFound)
}

func readPostedBody(w http.ResponseWriter, r *http.Request, maxBytes int64) creds {
	if ct := r.Header.Get("Content-Type"); ct != "" {
		if strings.ToLower(ct) != "application/json" {
			http.Error(w, "Content-Type header is not application/json", http.StatusBadRequest)
		}
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)

	credentials := new(creds)

	err := decoder.Decode(credentials)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {

		case errors.As(err, &syntaxError):
			fallthrough
		case errors.Is(err, io.ErrUnexpectedEOF):
			fallthrough

		case errors.As(err, &unmarshalTypeError):
			fallthrough

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fallthrough

		case errors.Is(err, io.EOF):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case err.Error() == "http: request body too large":
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)

		default:
			log.Print(err.Error())
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
	}
	return *credentials
}
