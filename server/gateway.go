package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

// A gateway is middleware that handles all requests that require authentication (except for initial login requests)
type gateway struct {
	handler http.Handler
}

func newGateway(h http.Handler) *gateway {
	return &gateway{handler: h}
}

// ServeHTTP validates the session and the token in the request cookies then passes the request to the handler.
func (g *gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("authorizing request", "url", r.URL.String())
	if strings.HasPrefix(r.URL.Path, homePath) && r.Method != http.MethodGet {
		http.Error(w, "need to GET the page", http.StatusMethodNotAllowed)
		return
	}

	sessionIDCookie, err := r.Cookie(cookieSession)
	if err != nil {
		slog.Debug("session cookie not found")
		http.Redirect(w, r, "/homegym/login/", http.StatusUnauthorized)
		return
	}

	tokenCookie, err := r.Cookie(cookieToken)
	if err != nil {
		slog.Debug("token not found")
		http.Redirect(w, r, "/homegym/login/", http.StatusUnauthorized)
		return
	}
	sessionID := strings.TrimSpace(sessionIDCookie.Value)
	token := strings.TrimSpace(tokenCookie.Value)

	validToken, err := authorizer.ValidateToken(token, sessionID)
	if err != nil {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	slog.Info("authorized ", "url", r.URL.String())

	// redirect to home page and set internal routing cookie
	for _, path := range internalRoutes {
		if strings.HasPrefix(r.URL.Path, path) {
			routingCookie := http.Cookie{
				Name:     cookieRoute,
				Value:    r.URL.Path,
				Secure:   isSafe,
				HttpOnly: isSafe,
				Path:     "/homegym/",
				SameSite: samesite,
			}
			http.SetCookie(w, &routingCookie)
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
		slog.Error(err.Error())
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, usernameKey, claims.Audience[0])
	ctx = context.WithValue(ctx, roleKey, claims.Role)

	r = r.Clone(ctx)

	g.handler.ServeHTTP(w, r)
}
