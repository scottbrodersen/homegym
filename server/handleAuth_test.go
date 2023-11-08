package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/scottbrodersen/homegym/auth"
	"github.com/scottbrodersen/homegym/workoutlog"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

var (
	testSessionID = "testSessionID"
	testPassword  = "testpassword"
	testUserName  = "testUserName"
	testEmail     = "testemail@example.com"
)

var testTokenCookie http.Cookie = http.Cookie{
	Name:     cookieToken,
	Value:    testToken,
	Path:     "/homegym/",
	SameSite: samesite,
	MaxAge:   authorizer.SessionTTL(),
}
var testSessionCookie http.Cookie = http.Cookie{
	Name:     cookieSession,
	Value:    testSessionID,
	Path:     "/homegym/",
	SameSite: samesite,
	MaxAge:   authorizer.SessionTTL(),
}

func TestHandleAuth(t *testing.T) {

	Convey("Given an authorizer", t, func() {
		mockAuth := NewMockAuthorizer()
		authorizer = mockAuth
		mockUserAdmin := newMockUserAdmin()

		Convey("When we receive a login request with valid credentials", func() {
			mockAuth.On("IssueToken", mock.Anything, mock.Anything).Return(&testToken, &testSessionID, nil)

			url := fmt.Sprintf("/homegym/login/?username=%s&password=%s", testUserName, testPassword)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()

			HandleLogin(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusFound)
			setCookie := w.Result().Header.Values("Set-Cookie")
			So(setCookie, ShouldNotBeEmpty)

			foundTokenHeader := false
			foundSessionHeader := false

			for _, v := range setCookie {
				So(strings.Contains(v, "Path=/homegym/"), ShouldBeTrue)
				ss := samesiteString()
				So(strings.Contains(v, fmt.Sprintf("SameSite=%s", ss)), ShouldBeTrue)

				if strings.Contains(v, fmt.Sprintf("%s=%s", cookieToken, testToken)) {
					foundTokenHeader = true
				} else if strings.Contains(v, fmt.Sprintf("%s=%s", cookieSession, testSessionID)) {
					foundSessionHeader = true
				}
			}

			So(foundSessionHeader, ShouldBeTrue)
			So(foundTokenHeader, ShouldBeTrue)
		})

		Convey("When we receive a request with the incorrect method", func() {
			url := fmt.Sprintf("/homegym/login/?username=%s&password=%s", testUserName, testPassword)
			req := httptest.NewRequest(http.MethodPut, url, nil)

			w := httptest.NewRecorder()

			HandleLogin(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusMethodNotAllowed)
			So(w.Result().Header.Values("Set-Cookie"), ShouldBeEmpty)
		})

		Convey("When we recieve a request with invalid credentials", func() {
			mockAuth.On("IssueToken", mock.Anything, mock.Anything).Return(nil, nil, auth.ErrUnauthorized)

			url := fmt.Sprintf("/homegym/login/?username=%s&password=%s", testUserName, testPassword)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()

			HandleLogin(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("When we receive a valid signup request", func() {
			workoutlog.FrontDesk = mockUserAdmin

			mockUserAdmin.On("NewUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&workoutlog.User{}, nil)

			url := fmt.Sprintf("/homegym/signup/?username=%s&password=%s&email=%s", testUserName, testPassword, testEmail)
			req := httptest.NewRequest(http.MethodPost, url, nil)

			w := httptest.NewRecorder()

			HandleSignup(w, req)

			So(w.Result().StatusCode, ShouldEqual, http.StatusFound)
			setCookie := w.Result().Header.Values("Set-Cookie")
			So(setCookie, ShouldNotBeEmpty)
			So(strings.Contains(setCookie[0], fmt.Sprintf("%s=%s", cookieUsername, testUserName)), ShouldBeTrue)
			So(strings.Contains(setCookie[0], "Path=/homegym/login/"), ShouldBeTrue)
			ss := samesiteString()
			So(strings.Contains(setCookie[0], fmt.Sprintf("SameSite=%s", ss)), ShouldBeTrue)
		})
	})
}
