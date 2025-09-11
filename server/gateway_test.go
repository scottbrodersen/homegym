package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/scottbrodersen/homegym/auth"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

var (
	testToken = "testToken"
)

var stdClaims jwt.RegisteredClaims = jwt.RegisteredClaims{Audience: []string{testUserName}}
var gymClaims auth.GymClaims = auth.GymClaims{Role: string(auth.User)}
var testClaims auth.Claims = auth.Claims{
	RegisteredClaims: stdClaims,
	GymClaims:        gymClaims,
}

func TestGateway(t *testing.T) {
	gw := gateway{handler: newTestHandler()}

	Convey("When a request has no session cookie", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/homegym/home/", nil)
		req.AddCookie(&testTokenCookie)

		w := httptest.NewRecorder()

		gw.ServeHTTP(w, req)
		So(w.Result().StatusCode, ShouldEqual, http.StatusUnauthorized)
		//so(w.Result().)
	})

	Convey("When a request has no token cookie", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/homegym/api/foo", nil)
		req.AddCookie(&testSessionCookie)

		w := httptest.NewRecorder()

		gw.ServeHTTP(w, req)
		So(w.Result().StatusCode, ShouldEqual, http.StatusUnauthorized)
	})
	Convey("When a request for the home page is a POST", t, func() {
		req := httptest.NewRequest(http.MethodPost, "/homegym/home/", nil)
		req.AddCookie(&testSessionCookie)
		req.AddCookie(&testTokenCookie)

		w := httptest.NewRecorder()

		gw.ServeHTTP(w, req)
		So(w.Result().StatusCode, ShouldEqual, http.StatusMethodNotAllowed)
	})

	Convey("Given an authorizer", t, func() {
		mockAuth := NewMockAuthorizer()
		authorizer = mockAuth
		Convey("When we request the home page using valid access tokens", func() {
			mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(&testToken, nil)

			req := httptest.NewRequest(http.MethodGet, "/homegym/home/", nil)
			req.AddCookie(&testTokenCookie)
			req.AddCookie(&testSessionCookie)

			w := httptest.NewRecorder()

			gw.ServeHTTP(w, req)
			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)
			setCookie := w.Result().Header.Get("Set-Cookie")
			So(setCookie, ShouldNotBeEmpty)
			So(strings.Contains(setCookie, fmt.Sprintf("%s=%s", cookieToken, testToken)), ShouldBeTrue)
			So(strings.Contains(setCookie, "Path=/homegym/"), ShouldBeTrue)
			ss := samesiteString()
			So(strings.Contains(setCookie, fmt.Sprintf("SameSite=%s", ss)), ShouldBeTrue)
		})

		Convey("When a request is not authorized", func() {
			mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test error"))

			req := httptest.NewRequest(http.MethodGet, "/homegym/home/", nil)
			req.AddCookie(&testTokenCookie)
			req.AddCookie(&testSessionCookie)

			w := httptest.NewRecorder()

			gw.ServeHTTP(w, req)
			So(w.Result().StatusCode, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("When we make an api call using valid access tokens", func() {
			mockAuth.On("ValidateToken", mock.Anything, mock.Anything).Return(&testToken, nil)
			mockAuth.On("TokenClaims", mock.Anything).Return(testClaims, nil)

			req := httptest.NewRequest(http.MethodGet, "/homegym/api/foo", nil)
			req.AddCookie(&testTokenCookie)
			req.AddCookie(&testSessionCookie)

			w := httptest.NewRecorder()

			gw.ServeHTTP(w, req)
			So(w.Result().StatusCode, ShouldEqual, http.StatusOK)

			body, _ := io.ReadAll(w.Result().Body)
			w.Result().Body.Close()

			expectedReqCtx := expectedCtx{}
			json.Unmarshal(body, &expectedReqCtx)

			So(expectedReqCtx.Role, ShouldEqual, testClaims.Role)
			So(expectedReqCtx.User, ShouldEqual, testClaims.Audience[0])
		})
	})

}

type testHandler string

// returns the request context as json
func (t testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := expectedCtx{
		Role: GymContextValue(req.Context(), roleKey),
		User: GymContextValue(req.Context(), usernameKey),
	}
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(ctx)

	w.Write(body)
}

func newTestHandler() testHandler {
	var th testHandler = "foo"
	return th
}

type expectedCtx struct {
	Role string
	User string
}
