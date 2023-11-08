package server

import (
	"context"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type gateway struct {
	handler http.Handler
}

// front end internal routing urls
var redirectToHome []string = []string{
	eventPath,
	activitiesPath,
	exerciseTypesPath,
}

func newGateway(h http.Handler) *gateway {
	return &gateway{handler: h}
}

func (g *gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.SetLevel(log.DebugLevel)
	log.Debug("authorizing request for ", r.URL.String())
	if strings.HasPrefix(r.URL.Path, homePath) && r.Method != http.MethodGet {
		http.Error(w, "need to GET the page", http.StatusMethodNotAllowed)
		return
	}

	sessionIDCookie, err := r.Cookie(cookieSession)
	if err != nil {
		if r.URL.Path == homePath {
			http.Redirect(w, r, "/homegym/login/", http.StatusFound)
			return
		}
		http.Error(w, `{"message": "Log in again please"}`, http.StatusUnauthorized)
		return
	}

	tokenCookie, err := r.Cookie(cookieToken)
	if err != nil {
		log.Debug("token not found, returning 401")
		if r.URL.Path == homePath {
			http.Redirect(w, r, "/homegym/login/", http.StatusFound)
			return
		}
		http.Error(w, `{"message": "Log in again please"}`, http.StatusUnauthorized)
		return
	}
	sessionID := strings.TrimSpace(sessionIDCookie.Value)
	token := strings.TrimSpace(tokenCookie.Value)

	validToken, err := authorizer.ValidateToken(token, sessionID)
	if err != nil {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	log.Debug("authorized for ", r.URL.String())

	// redirect to home page
	for _, path := range redirectToHome {
		if strings.HasPrefix(r.URL.Path, path) {
			http.Redirect(w, r, homePath, http.StatusFound)
			return
		}
	}

	if strings.HasPrefix(r.URL.Path, homePath) {
		if r.URL.Path == homePath {
			r.URL.Path += "dist/"
		} else if strings.HasPrefix(r.URL.Path, homePath+"assets/") {
			r.URL.Path = strings.Replace(r.URL.Path, "/assets/", "/dist/assets/", 1)
		}

		http.SetCookie(w, tokenCookieMaker(validToken))
		g.handler.ServeHTTP(w, r)

		return
	}

	// this is an api call
	// add info to context for later authz
	claims, err := authorizer.TokenClaims(token)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, `{"message": "something went wrong"}`, http.StatusInternalServerError)
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, usernameKey, claims.Audience[0])
	ctx = context.WithValue(ctx, roleKey, claims.Role)

	r = r.Clone(ctx)

	g.handler.ServeHTTP(w, r)
}
