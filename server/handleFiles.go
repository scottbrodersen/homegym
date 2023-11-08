package server

import (
	"io/fs"
	"net/http"
	"strings"
)

func GymFileServer(embeddedFS fs.FS) http.HandlerFunc {
	fs := http.FileServer(http.FS(embeddedFS))
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.RequestURI, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasPrefix(r.RequestURI, "index") && strings.HasSuffix(r.RequestURI, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		}
		fs.ServeHTTP(w, r)
	}
}
