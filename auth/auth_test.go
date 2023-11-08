package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/scottbrodersen/homegym/dal"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

var (
	testPassword    = "testpassword0123456789"
	shortPassword   = "short"
	longPassword    = "longpasswordlongpasswordlongpasswordlongpasswordlongpasswordlongpasswordlongpassword" //converts to 84-byte byte slice
	passwordVersion = "v1"
)

var a Authorizer = Authorizer{}

func TestPasswords(t *testing.T) {
	// Test helps ensure latest version is up to date
	Convey("When we get the latest version", t, func() {
		latest := LatestVersion()
		So(latest, ShouldNotBeEmpty)
		So(latest, ShouldEqual, "v1")
		Convey("And when we get the hash and validate functions for the latest versionJ", func() {
			v, okv := validateFunctions[latest]
			h, okh := hashFunctions[latest]
			So(okv, ShouldBeTrue)
			So(okh, ShouldBeTrue)
			So(v, ShouldEqual, validatePasswordVersion1)
			So(h, ShouldEqual, hashPasswordVersion1)
		})
	})
	for _, version := range [1]string{"v1"} {
		Convey("Given a v1 password utility", t, func() {
			pwdUtil := PasswordUtil{Version: version}
			Convey("When we hash a short password using a versioned hasher", func() {
				hash, err := pwdUtil.Hash(shortPassword)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, errPwdShort)
				So(hash, ShouldBeEmpty)
			})
			Convey("When we hash a long password using v1 hasher", func() {
				hash, err := pwdUtil.Hash(longPassword)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, errPwdLong)
				So(hash, ShouldBeEmpty)
			})
			Convey("When we hash a valid password using v1 hasher", func() {
				hash, err := pwdUtil.Hash(testPassword)
				So(err, ShouldBeNil)
				So(hash, ShouldNotBeEmpty)
				Convey("And when we validate the password using the v1 validator", func() {
					ok := pwdUtil.ValidatePassword(testPassword, hash)
					So(ok, ShouldBeNil)
				})
				Convey("And when we validate the wrong password using the v1 validator", func() {
					err := pwdUtil.ValidatePassword("bad", hash)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, errPwdWrong)
				})
			})
		})
	}
}

var (
	testUserName  = "testUsername"
	testEmail     = "testemail@example.co"
	testSessionID = "test-session-id"
)

var testRole Role = "TestRole"

func TestTokens(t *testing.T) {
	pwdUtil := PasswordUtil{Version: passwordVersion}
	testPwdHash, _ := pwdUtil.Hash(testPassword)
	testTokenKey, _ := generateRsaPrivate(2048)
	testTokenKeyBytes := rsaKeyToByteSlice(testTokenKey)
	testKeyID := "test-key-id"

	var testToken *string
	var err error
	var sessionID *string
	testExpiry := time.Now().Add(time.Second * 3000).Unix()
	testSession := &Session{
		UserID:    testUserName,
		SessionID: testSessionID,
		Expiry:    testExpiry,
	}
	Convey("When we issue a token", t, func() {
		db := dal.NewMockDal()
		dal.DB = db
		roleStr := string(testRole)
		db.On("ReadUser", mock.Anything).Return(&testEmail, &testPwdHash, &passwordVersion, &roleStr, nil)
		db.On("GetKeys", mock.Anything).Return(map[string][]byte{testKeyID: testTokenKeyBytes}, map[string][]byte{}, nil)
		db.On("AddSession", mock.Anything, mock.Anything).Return(nil)

		beforeToken := time.Now().Add(time.Second * time.Duration(tokenLifeInSeconds))
		time.Sleep(time.Second * time.Duration(1))
		testToken, sessionID, err = a.IssueToken(testUserName, testPassword)
		afterToken := time.Now().Add(time.Second * time.Duration(tokenLifeInSeconds))
		time.Sleep(time.Second * time.Duration(1))

		So(err, ShouldBeNil)
		So(testToken, ShouldNotBeNil)
		So(sessionID, ShouldNotBeNil)
		claimPart, _ := base64.StdEncoding.DecodeString(strings.Split(*testToken, ".")[1])
		claims := Claims{}
		err = json.Unmarshal(claimPart, &claims)
		if err != nil {
			log.Fatal(err)
		}

		So(claims.Valid(), ShouldBeNil)
		So(claims.Audience, ShouldResemble, jwt.ClaimStrings{testUserName})
		So(claims.Role, ShouldEqual, string(testRole))
		So(claims.ExpiresAt.Time.After(beforeToken), ShouldBeTrue)
		So(claims.ExpiresAt.Time.Before(afterToken), ShouldBeTrue)

		Convey("And when we validate the token against a mismatched token", func() {
			sessionID := "random-id"
			expiry := testExpiry

			db.On("GetSession", mock.Anything, mock.Anything).Return(&sessionID, &expiry, nil)
			db.On("GetKeys", mock.Anything).Return(map[string][]byte{testKeyID: testTokenKeyBytes}, map[string][]byte{}, nil)

			expectNil, err := a.ValidateToken(*testToken, sessionID)

			So(err, ShouldNotBeNil)
			So(expectNil, ShouldBeNil)
		})

		Convey("And when we validate the token with the correct session and there are two token keys", func() {
			db.On("GetSession", mock.Anything, mock.Anything).Return(&testSession.UserID, &testSession.Expiry, nil)

			newTokenKey, _ := generateRsaPrivate(2048)
			newTokenKeyBytes := rsaKeyToByteSlice(newTokenKey)
			newTokenKeyID := "newTokenKeyID"

			db.On("GetKeys", mock.Anything).Return(map[string][]byte{testKeyID: testTokenKeyBytes, newTokenKeyID: newTokenKeyBytes}, map[string][]byte{}, nil)

			// wait so that refreshed token has updated expiry
			time.Sleep(time.Second * time.Duration(1))

			refreshed, err := a.ValidateToken(*testToken, *sessionID)

			So(err, ShouldBeNil)
			So(refreshed, ShouldNotBeNil)

			refreshedClaimPart, _ := base64.StdEncoding.DecodeString(strings.Split(*refreshed, ".")[1])
			refreshedClaims := Claims{}
			_ = json.Unmarshal(refreshedClaimPart, &refreshedClaims)

			So(refreshedClaims.ExpiresAt.Time.After(afterToken), ShouldBeTrue)
		})
	})

	Convey("And when we validate a token that's signed with a different key", t, func() {
		db := dal.NewMockDal()
		dal.DB = db

		db.On("GetSession", mock.Anything, mock.Anything).Return(&testSession.UserID, &testSession.Expiry, nil)
		db.On("GetKeys", mock.Anything).Return(map[string][]byte{testKeyID: testTokenKeyBytes}, map[string][]byte{}, nil)

		newKey, _ := generateRsaPrivate(2048)
		fakeTokenClaims := Claims{
			jwt.RegisteredClaims{
				ExpiresAt: tokenExpiryTime(),
				Issuer:    issuer,
				Audience:  []string{testUserName},
			},
			GymClaims{Role: "some role"},
		}

		fakeToken := jwt.NewWithClaims(jwt.SigningMethodRS256, fakeTokenClaims)
		fakeTokenStr, _ := fakeToken.SignedString(newKey)
		fakeRefreshed, err := a.ValidateToken(fakeTokenStr, testSessionID)

		So(err, ShouldNotBeNil)
		So(fakeRefreshed, ShouldBeNil)
	})

	Convey("When we issue a token with a user that isn't in the database", t, func() {
		db := dal.NewMockDal()
		dal.DB = db
		db.On("ReadUser", mock.Anything).Return(nil, nil, nil, nil, fmt.Errorf("test error"))
		token, sessionID, err := a.IssueToken(testUserName, testPassword)

		So(err, ShouldNotBeNil)
		So(token, ShouldBeNil)
		So(sessionID, ShouldBeNil)
	})

	Convey("When we validate a token against an unauthenticated user ID", t, func() {
		db := dal.NewMockDal()
		dal.DB = db
		db.On("GetSession", mock.Anything, mock.Anything).Return(nil, nil, nil)
		expectNil, err := a.ValidateToken(*testToken, "some-session-ID")
		So(err, ShouldNotBeNil)
		So(expectNil, ShouldBeNil)
	})

	Convey("When we extract the claims from a token", t, func() {
		testClaims := Claims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now()),
				Issuer:    "test issuer",
				Audience:  []string{testUserName},
			},
			GymClaims{Role: "test role"},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, testClaims)
		testKey, _ := generateRsaPrivate(2048)
		tokenStr, _ := token.SignedString(testKey)
		claims, err := a.TokenClaims(tokenStr)
		So(err, ShouldBeNil)
		So(claims, ShouldResemble, testClaims)
	})
}

func TestScheduledTasks(t *testing.T) {
	Convey("When we schedule a session for deletion", t, func() {
		duration := time.Millisecond * time.Duration(500)
		db := dal.NewMockDal()
		dal.DB = db
		db.On("DeleteSession", mock.Anything).Return(nil)
		endSessionAfterDuration(testSessionID, duration)
		db.Mock.AssertNotCalled(t, "DeleteSession", mock.Anything)
		time.Sleep(time.Millisecond * 300)
		db.Mock.AssertNotCalled(t, "DeleteSession", mock.Anything)
		time.Sleep(time.Millisecond * 300)
		db.Mock.AssertCalled(t, "DeleteSession", mock.Anything)
	})

	Convey("When we schedule a key for deletion", t, func() {
		duration := time.Millisecond * time.Duration(500)
		db := dal.NewMockDal()
		dal.DB = db

		db.On("DeleteKey", mock.Anything, mock.Anything).Return(nil)

		deleteKeyAfterDuration("keyID", KeyTypes.Token, duration)

		db.Mock.AssertNotCalled(t, "DeleteKey", mock.Anything, mock.Anything)

		time.Sleep(time.Millisecond * 300)

		db.Mock.AssertNotCalled(t, "DeleteKey", mock.Anything, mock.Anything)

		time.Sleep(time.Millisecond * 300)

		db.Mock.AssertCalled(t, "DeleteKey", mock.Anything, mock.Anything)
	})

	Convey("When we initiate key rotation for the first time", t, func() {
		{
			db := dal.NewMockDal()
			dal.DB = db
			firstKey, _ := generateRsaPrivate(2048)
			firstKeyBytes := rsaKeyToByteSlice(firstKey)
			secondKey, _ := generateRsaPrivate(2048)
			secondKeyBytes := rsaKeyToByteSlice(secondKey)
			firstKeyID := "first-Key-ID"
			secondKeyID := "second-Key-ID"

			briefDelay := 1 // seconds until scheduled rotation time
			schedule := time.Now().Unix() + int64(briefDelay)
			rotTime := time.Unix(schedule, 0).UTC().Format(time.TimeOnly)

			err := SetKeyRotationTime(rotTime)
			So(err, ShouldBeNil)

			// GetKeys returns none when initiating key rotation and when immediately rotating keys
			db.On("GetKeys", mock.Anything).Return(map[string][]byte{}, map[string][]byte{}, nil).Twice()
			// RotateKeys is called when immediately rotating keys (when initiating key rotation) and when the scheduled rotation occurs
			db.On("RotateKeys", mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()
			// GetKeys is called again when the scheduled rotation occurs before deleting the old keys
			db.On("GetKeys", mock.Anything).Return(map[string][]byte{secondKeyID: secondKeyBytes}, map[string][]byte{firstKeyID: firstKeyBytes}, nil)
			// DeleteKey is called when deleting the old keys after the scheduled rotation
			db.On("DeleteKey", mock.Anything, mock.Anything).Return(nil)

			secondsToDelayKeyDeletion = 1

			// There are no keys in the database so initiating key rotation should rotate keys immediately then schedule the daily rotation
			err = InitiateKeyRotation(KeyTypes.Token) //returns after rotating key and scheduling the daily rotation

			// set the daily rotation time to occur at a much later time so it does not interfere with our test
			longDelay := 20
			reschedule := time.Now().Unix() + int64(longDelay)
			rotTime2 := time.Unix(reschedule, 0).UTC().Format(time.TimeOnly)
			_ = SetKeyRotationTime(rotTime2)

			// The immediate key rotation is done now
			So(err, ShouldBeNil)
			So(db.Mock.AssertNumberOfCalls(t, "GetKeys", 2), ShouldBeTrue)
			So(db.Mock.AssertNumberOfCalls(t, "RotateKeys", 1), ShouldBeTrue)
			So(db.Mock.AssertNumberOfCalls(t, "DeleteKey", 0), ShouldBeTrue)

			// Wait for the daily rotation to occur
			time.Sleep(time.Second * time.Duration(briefDelay+1))

			So(db.Mock.AssertNumberOfCalls(t, "RotateKeys", 2), ShouldBeTrue)
			So(db.Mock.AssertNumberOfCalls(t, "GetKeys", 3), ShouldBeTrue)
			So(db.Mock.AssertNumberOfCalls(t, "DeleteKey", 0), ShouldBeTrue)

			time.Sleep(time.Second * time.Duration(secondsToDelayKeyDeletion))
			So(db.Mock.AssertNumberOfCalls(t, "GetKeys", 3), ShouldBeTrue)
			So(db.Mock.AssertNumberOfCalls(t, "DeleteKey", 1), ShouldBeTrue)
			So(db.Mock.AssertNumberOfCalls(t, "RotateKeys", 2), ShouldBeTrue)
		}
	})
}
