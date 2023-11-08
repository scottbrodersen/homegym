package workoutlog

import (
	"testing"

	"github.com/scottbrodersen/homegym/auth"

	"github.com/scottbrodersen/homegym/dal"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

var (
	testUserName    = "testusername"
	testRoleStr     = "TestRole"
	testRole        = auth.Role(testRoleStr)
	testEmail       = "testemail@example.co"
	testPassword    = "testpassword0123456789"
	testHash        = testPassword + "hash"
	passwordVersion = "mock"
)

func testUser() User {
	return User{
		ID:             testUserName,
		Role:           testRole,
		Email:          testEmail,
		PwdHash:        testHash,
		PwdHashVersion: passwordVersion,
	}
}

func TestUsers(t *testing.T) {

	Convey("Given a dal client", t, func() {
		db := dal.NewMockDal()
		dal.DB = db
		Convey("When we add a user with correct attributes", func() {
			db.On("NewUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			user, err := FrontDesk.NewUser(testUserName, testRole, testEmail, testPassword)

			Convey("Then a user with the expected attributes is returned", func() {
				So(err, ShouldBeNil)
				So(user, ShouldNotBeNil)
				So(user.ID, ShouldEqual, testUserName)
				So(user.Role, ShouldEqual, testRole)
				So(user.PwdHash, ShouldNotBeEmpty)
				So(user.PwdHashVersion, ShouldEqual, auth.LatestVersion())
			})
		})

		Convey("When we Get the user", func() {
			testUser := testUser()
			db.On("ReadUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&testEmail, &testHash, &passwordVersion, &testRoleStr, nil)
			user, err := FrontDesk.GetUser(testUserName)
			Convey("The user object is returned", func() {
				So(err, ShouldBeNil)
				So(user, ShouldResemble, &testUser)
			})
		})
		Convey("When we update the user with the same password", func() {
			user := testUser()
			db.On("ReadUser", mock.Anything).Return(&testEmail, &testHash, &passwordVersion, &testRoleStr, nil)
			db.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			err := user.UpdateUserPassword(testPassword)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When we update the user with a new password", func() {
			user := testUser()
			hash := user.PwdHash
			db.On("ReadUser", mock.Anything).Return(&testEmail, &testHash, &passwordVersion, &testRoleStr, nil)
			db.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			err := user.UpdateUserPassword("newpassword")
			Convey("Then the user is updated as expected", func() {
				So(err, ShouldBeNil)
				So(hash, ShouldNotEqual, user.PwdHash)
			})
		})
		Convey("When we update the user details", func() {
			user := testUser()
			user.Email = "new@example.com"
			db.On("ReadUser", mock.Anything).Return(&testEmail, &testHash, &passwordVersion, &testRoleStr, nil)
			db.On("UpdateUserProfile", mock.Anything, mock.Anything).Return(nil)
			err := user.UpdateUserProfile()
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})

	})
}
