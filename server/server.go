// Package server implements the HTTP server for Homegym.

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/scottbrodersen/homegym/auth"
	"github.com/scottbrodersen/homegym/dal"
	"github.com/scottbrodersen/homegym/server/public"
	"github.com/scottbrodersen/homegym/server/secured"
)

const (
	usernameKey         GymContextKey = "username"
	roleKey             GymContextKey = "role"
	cookieToken                       = "token"
	cookieSession                     = "session"
	cookieUsername                    = "username"
	cookieRoute                       = "followroute"
	allowedCorsOrigin                 = "http://127.0.0.1:3000"
	allowedCorsMethods                = "*"
	allowedCorsHeaders                = "Set-Cookie"
	internalServerError               = `{"message":"something went wrong"}`
)

var samesite http.SameSite = http.SameSiteLaxMode
var isSafe bool = false

type requestAuthorizer interface {
	IssueToken(username string, pwd string) (*string, *string, error)
	ValidateToken(tokenString, sessionID string) (*string, error)
	TokenClaims(tokenString string) (auth.Claims, error)
	TokenTTL() int
	SessionTTL() int
}

const homePath string = "/homegym/home/"

// frontend app internal routes
var internalRoutes []string = []string{
	"/homegym/event/",
	"/homegym/activities/",
	"/homegym/exercises/",
	"/homegym/programs/",
}

var secureMux *http.ServeMux
var publicMux *http.ServeMux
var secureGateway http.Handler
var authorizer requestAuthorizer = auth.NewAuthorizer()

// We use middleware as a bridge between a public mux and secured mux.
// The public mux routes top-level paths to middleware that authenticates the request (a gateway)
// Once authenticated, the middleware passes the request to the secured mux.
// Routes to the login page and signup page are unauthenticated.
func init() {
	// Routes accessible after authentication by secureGateway
	secureMux = http.NewServeMux()

	secureMux.HandleFunc("/homegym/api/activities/", ActivitiesApi)
	secureMux.HandleFunc("/homegym/api/exercises/", ExerciseTypesApi)
	secureMux.HandleFunc("/homegym/api/events/", EventsApi)
	secureMux.HandleFunc("/homegym/api/dailystats/", DailyStatsApi)
	secureMux.Handle("/homegym/home/dist/", http.StripPrefix("/homegym/home", GymFileServer(secured.SecuredEFS)))

	// middleware that authenticates before relaying to secure mux
	secureGateway = newGateway(secureMux)

	// handler for requests to http server
	publicMux = http.NewServeMux()

	// least specific path -- resolves to static public assets
	publicMux.Handle("/homegym/", http.StripPrefix("/homegym", GymFileServer(public.HtmlEFS)))
	log.SetLevel(log.DebugLevel)

	// first-level subtrees -- routed to secureGateway for authentication
	publicMux.Handle("/homegym/api/", secureGateway)
	publicMux.Handle("/homegym/home/", secureGateway)

	for _, internalRoute := range internalRoutes {
		publicMux.Handle(internalRoute, secureGateway)
	}

	// specific paths for initial authentication requests
	publicMux.HandleFunc("/homegym/login", HandleLogin)
	publicMux.HandleFunc("/homegym/signup", HandleSignup)
}

func standardHeaders(header *http.Header) {
	//TODO: add csp
	header.Add("content-type", "application/json")
	//header.Add("Access-Control-Allow-Origin", allowedCorsOrigin)
	//header.Add("Cross-Origin-Resource-Policy", "same-origin")
	header.Add("Access-Control-Allow-Methods", allowedCorsMethods)
	header.Add("Access-Control-Allow-Headers", allowedCorsHeaders)
	header.Add("Access-Control-Allow-Credentials", "true")
}

func StartUnsafe(shutdown shutdownAction, port int) {
	shutdown(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), NewRequestLogger(publicMux)))
}

func StartSafe(shutdown shutdownAction) {
	//shutdown(http.ListenAndServeTLS())
	isSafe = true
}

type shutdownAction func(err error)

var DefaultShutdown shutdownAction = func(err error) {
	// do some cleanup, and then log.fatal
	dal.DB.Destroy()
	log.Fatal(err)
}

type GymContextKey string

func GymContextValue(ctx context.Context, k GymContextKey) string {
	if v := ctx.Value(k); v != nil {
		return v.(string)
	}
	return ""
}

func tokenCookieMaker(token *string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieToken,
		Value:    *token,
		Secure:   isSafe,
		HttpOnly: isSafe,
		Path:     "/homegym/",
		SameSite: samesite,
		MaxAge:   authorizer.TokenTTL(),
	}
}

func whoIsIt(ctx context.Context) (*string, *string, error) {
	username := GymContextValue(ctx, usernameKey)
	if username == "" {
		return nil, nil, fmt.Errorf("no username")
	}

	role := GymContextValue(ctx, roleKey)
	if role == "" {
		return nil, nil, fmt.Errorf("no role")
	}

	return &username, &role, nil
}

func jsonSafeError(e error) string {
	replacement := ": "

	var replacer = strings.NewReplacer(
		"\r\n", replacement,
		"\r", replacement,
		"\n", replacement,
		"\v", replacement,
		"\f", replacement,
		"\u0085", replacement,
		"\u2028", replacement,
		"\u2029", replacement,
	)
	r := replacer.Replace(e.Error())
	return r
}
