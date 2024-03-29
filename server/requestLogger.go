package server

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

type RequestLogger struct {
	handler http.Handler
	log     *logrus.Logger
}

func NewRequestLogger(h http.Handler) *RequestLogger {
	logger := logrus.New()
	return &RequestLogger{handler: h, log: logger}
}

// ServeHTTP logs request details then passes the request to its handler
func (rl *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fields := map[string]interface{}{"requested": r.RequestURI, "cookies": r.Header.Values("Cookie"), "remote": r.RemoteAddr}
	scrubPII(fields, scrubPassword, scrubCookies)
	rl.log.WithFields(logrus.Fields(fields)).Info("Handling request")

	rl.handler.ServeHTTP(w, r)
}

func scrubPII(fields map[string]interface{}, scrubbers ...scrubber) {
	for k, v := range fields {
		switch val := v.(type) {
		case string:
			for _, scrubber := range scrubbers {
				fields[k] = scrubber(val)
			}
		case []string:
			newvals := []string{}
			for _, value := range val {
				for _, scrubber := range scrubbers {
					value = scrubber(value)
				}
				newvals = append(newvals, value)
			}
			fields[k] = newvals
		}
	}

}

type scrubber func(string) string

var scrubPassword scrubber = func(s string) string {
	p := regexp.MustCompile(`(?m)^(?:.*[&?])((password|pwd)\s*[=:]\s*)([^;&]+)[;&]?$`)
	return string(p.ReplaceAll([]byte(s), []byte("$1{redacted}")))
}

var scrubCookies scrubber = func(s string) string {
	l := regexp.MustCompile(`((token|session)=\s*[^;]{5})[^;]+(;?)`)

	return fmt.Sprint(string(l.ReplaceAll([]byte(s), []byte("$1...$3"))))
}
