package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scottbrodersen/homegym/dal"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	UserID    string
	SessionID string
	Expiry    int64
}

type Role string

type CryptoKey struct {
	Key     *rsa.PrivateKey
	Retired bool
}

var Admin Role = "admin"
var User Role = "healthnut"

// Time of daily key rotation in format "hh:mm:ss".
var keyRotationTime = "23:59:00"
var secondsToDelayKeyDeletion = 300 // 5 minutes
const tokenLifeInSeconds = 3000     // 5 mins -- extended for testing
const sessionLifeInMinutes = 1440   // 24 hours

// Set the time of the daily key rotation.
// Default time is 23:59:00
func SetKeyRotationTime(rotationTime string) error {
	rx := regexp.MustCompile("^[0-9]{2}:[0-9]{2}:[0-9]{2}$")
	if !rx.Match([]byte(rotationTime)) {
		err := fmt.Errorf("invalid time, using default: %s", keyRotationTime)
		log.WithError(err).Error("key rotation time not in hh:mm:ss format")
		log.Info("using default rotation time.")
		return err
	}
	log.Info("setting key rotation time to: ", rotationTime)
	keyRotationTime = rotationTime
	return nil
}

type KeyUsage struct {
	Token  string
	Pepper string
}

var KeyTypes KeyUsage = KeyUsage{
	Token:  "token",
	Pepper: "pepper",
}

// Rotates keys immediately if no usable key is found.
// Starts the daily key rotation
func InitiateKeyRotation(keyType string) error {
	// get current key and if there is none create one now
	if keyType == KeyTypes.Token {
		signingKey, err := getSigningKey()
		if err != nil {
			log.WithError(err).Error("failed to get signing key -- ignoring")
		}
		if signingKey == nil {
			rotateKey(KeyTypes.Token)
		}
	}
	// initiate daily rotation
	log.Info("initiating daily key rotation.")
	dailyKeyRotation(keyType)

	return nil
}

// Gets all token keys and returns the (first) one with UserMeta of 0
func getSigningKey() (*rsa.PrivateKey, error) {
	active, _, err := dal.DB.GetKeys(KeyTypes.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to read keys: %w", err)
	}

	for _, key := range active {

		keyRSA, err := byteSliceToRSAKey(key)
		if err != nil {
			return nil, err
		}

		return keyRSA, nil
	}
	return nil, nil
}

// Calculates the time until the next scheduled key rotation
// Schedules the next key rotation.
// Rinse and repeat.
func dailyKeyRotation(keyType string) {
	now := time.Now().UTC()
	y, m, d := now.Date()

	mPrefix := ""
	dPrefix := ""

	if int(m) < 10 {
		mPrefix = "0"
	}
	if d < 10 {
		dPrefix = "0"
	}
	// 	RFC3339     = "2006-01-02T15:04:05Z07:00"
	nextRotationTime, err := time.Parse(time.RFC3339, fmt.Sprintf("%d-%s%d-%s%dT%sZ", y, mPrefix, m, dPrefix, d, keyRotationTime))
	if err != nil {
		log.WithError(err).Error("failed to calculate time of next key rotation")
	}

	duration := time.Second * time.Duration(nextRotationTime.Unix()-now.Unix())
	log.Info(fmt.Sprintf("Scheduling key rotation to occur in %s", duration.String()))

	time.AfterFunc(duration, func() {
		log.Info("rotating key on schedule: ", keyType)
		err := rotateKey(keyType)
		if err != nil {
			log.WithError(err).Error("failed to rotate key")
		}
		dailyKeyRotation(keyType)
	})
}

func rotateKey(keyType string) error {
	log.Info("rotating key")
	keyID := uuid.NewString()
	newKey, err := generateRsaPrivate(2048)
	if err != nil {
		return fmt.Errorf("could not create new key: %w", err)
	}

	newKeyByte := rsaKeyToByteSlice(newKey)

	err = dal.DB.RotateKeys(newKeyByte, keyID, keyType)
	if err != nil {
		return err
	}

	_, existingKeys, err := dal.DB.GetKeys(keyType)
	if err != nil {
		return fmt.Errorf("failed to get keys during rotation: %w", err)
	}
	for id := range existingKeys {
		deleteKeyAfterDuration(id, keyType, time.Second*time.Duration(secondsToDelayKeyDeletion))
	}

	return nil
}

func deleteKeyAfterDuration(keyID, usage string, duration time.Duration) {
	log.Info("scheduling deletion of key: ", usage)

	time.AfterFunc(duration, func() {
		log.Info("deleting key on schedule: ", usage)

		err := dal.DB.DeleteKey(keyID, usage)
		if err != nil {
			log.WithError(err).Errorf("failed to delete %s key", usage)
		}
	})
}

func generateRsaPrivate(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func rsaKeyToByteSlice(k *rsa.PrivateKey) []byte {
	keyX509 := x509.MarshalPKCS1PrivateKey(k)
	keyStr := base64.StdEncoding.EncodeToString(keyX509)
	return []byte(keyStr)
}

func byteSliceToRSAKey(byteKey []byte) (*rsa.PrivateKey, error) {
	keyX509, err := base64.StdEncoding.DecodeString(string(byteKey))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key: %w", err)
	}
	keyRSA, err := x509.ParsePKCS1PrivateKey(keyX509)
	if err != nil {
		return nil, fmt.Errorf("failed to parse x509 key: %w", err)
	}
	return keyRSA, nil
}

/* PASSWORDS */
type hashPasswordVersion func(string) (string, error)
type validatePasswordVersion func(string, string) error
type PasswordUtil struct {
	Version string
}

var (
	errPwdShort = "password too short"
	errPwdLong  = "password too long"
	errPwdWrong = "bad credentials"
)
var hashFunctions map[string]hashPasswordVersion = map[string]hashPasswordVersion{
	"v1":   hashPasswordVersion1,
	"mock": hashPasswordMock,
}
var validateFunctions map[string]validatePasswordVersion = map[string]validatePasswordVersion{
	"v1":   validatePasswordVersion1,
	"mock": validatePasswordMock,
}

func LatestVersion() string {
	return "v1"
}

func (p *PasswordUtil) Hash(pwd string) (string, error) {
	if hashFunc, ok := hashFunctions[p.Version]; ok {
		return hashFunc(pwd)
	}
	return "", fmt.Errorf("invalid version")
}

func (p *PasswordUtil) ValidatePassword(pwd string, hash string) error {
	if validateFunc, ok := validateFunctions[p.Version]; ok {
		return validateFunc(pwd, hash)
	}
	return fmt.Errorf("invalid version")
}

func hashPasswordVersion1(pwd string) (string, error) {
	minPwdLength := 10
	maxPwdLength := 72 //bytes
	if len(pwd) < minPwdLength {
		return "", fmt.Errorf(errPwdShort)
	}
	pwdBytes := []byte(pwd)
	if len(pwdBytes) > maxPwdLength {
		return "", fmt.Errorf(errPwdLong)
	}
	hash, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func validatePasswordVersion1(pwd string, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)); err != nil {
		return fmt.Errorf(errPwdWrong)
	}
	return nil
}

func hashPasswordMock(pwd string) (string, error) {
	if pwd != "" {
		return fmt.Sprintf("%shash", pwd), nil
	}
	return "", fmt.Errorf("mock error")
}

func validatePasswordMock(pwd, hash string) error {
	if hash != pwd+"hash" {
		return fmt.Errorf("mock error")
	}
	return nil
}

/* TOKENS */

type Authorizer struct {
	tokenTTL   int
	sessionTTL int
}

func (a Authorizer) TokenTTL() int {
	return a.tokenTTL
}

func (a Authorizer) SessionTTL() int {
	return a.sessionTTL * 60
}

func NewAuthorizer() Authorizer {
	return Authorizer{
		tokenTTL:   tokenLifeInSeconds,
		sessionTTL: sessionLifeInMinutes,
	}
}

const (
	issuer = "homegym"
)

var ErrUnauthorized error = errors.New("authentication failed")

type Claims struct {
	jwt.RegisteredClaims
	GymClaims
}
type GymClaims struct {
	Role string
}

func (c *Claims) Valid() error {
	if c.Issuer != issuer {
		log.Debugf("token claim contained invalid issuer: %s", c.Issuer)
		return fmt.Errorf("invalid issuer")
	}
	return nil
}

func tokenExpiryTime() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(tokenLifeInSeconds)))
}

// Authenticates user credentials.
// Creates a token.
// Stores the session.
// Returns the token string, sessionID, error
func (a Authorizer) IssueToken(username string, pwd string) (*string, *string, error) {
	_, pwdHash, pwdHashVersion, role, err := dal.DB.ReadUser(username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}
	pwdUtil := PasswordUtil{Version: *pwdHashVersion}
	if err = pwdUtil.ValidatePassword(pwd, *pwdHash); err != nil {
		return nil, nil, ErrUnauthorized
	}

	signingKey, err := getSigningKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get signing key: %w", err)
	}
	if signingKey == nil {
		return nil, nil, fmt.Errorf("no signing key")
	}

	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: tokenExpiryTime(),
			Issuer:    issuer,
			Audience:  []string{username},
		},
		GymClaims{Role: string(*role)},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create token: %w", err)
	}

	sessionID := uuid.NewString()
	err = dal.DB.AddSession(username, sessionID, a.sessionTTL)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create user session: %w", err)
	}
	log.Info("Session created: ", sessionID)
	endSessionAfterDuration(sessionID, time.Hour*time.Duration(24))

	return &tokenStr, &sessionID, nil
}

// Checks if the user's session is valid.
// Checks if the user's token is valid.
// Returns a new token with updated expiry.
func (a Authorizer) ValidateToken(tokenString, sessionID string) (*string, error) {

	username, expiry, err := dal.DB.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("could not authenticate user: %w", err)
	}
	if username == nil || expiry == nil || time.Now().Unix() > *expiry {
		return nil, fmt.Errorf("session expired")
	}

	claims := &Claims{}
	parser := jwt.NewParser(jwt.WithAudience(*username), jwt.WithIssuer(issuer), jwt.WithValidMethods([]string{"RS256"}))

	active, _, err := dal.DB.GetKeys(KeyTypes.Token)

	var token *jwt.Token = nil
	var parseErr error = nil
	var rsaKey *rsa.PrivateKey = nil

	for _, k := range active {
		rsaKey, err = byteSliceToRSAKey(k)
		if err != nil {
			return nil, fmt.Errorf("could not parse key: %w", err)
		}

		break
	}

	publicKey := rsaKey.PublicKey

	token, parseErr = parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return &publicKey, nil })
	if parseErr != nil {
		return nil, fmt.Errorf("could not parse token: %w", parseErr)
	}

	var ok bool
	claims, ok = token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("could not read claims: %w", err)
	}

	exp, _ := claims.GetExpirationTime()
	if exp.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	// refresh token
	refreshedClaims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: tokenExpiryTime(),
			Issuer:    issuer,
			Audience:  claims.Audience,
		},
		claims.GymClaims,
	}

	signingKey, err := getSigningKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	if signingKey == nil {
		return nil, fmt.Errorf("no signing key found")
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshedClaims)
	tokenStr, err := newToken.SignedString(signingKey)
	if err != nil {
		return nil, fmt.Errorf("could not create token: %w", err)
	}

	return &tokenStr, nil
}

func (a Authorizer) TokenClaims(tokenString string) (Claims, error) {
	parser := jwt.NewParser()
	var c Claims = Claims{}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("token does not have 3 parts")
	}

	tokenByte, err := parser.DecodeSegment(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("failed to parse token: %w", err)
	}
	err = json.Unmarshal(tokenByte, &c)
	if err != nil {
		return Claims{}, fmt.Errorf("failed to unmarshal claims: %w", err)

	}
	return c, nil

}

func endSessionAfterDuration(sessionID string, duration time.Duration) {
	log.Info("Scheduling deletion of session: ", duration)
	time.AfterFunc(duration, func() {
		err := dal.DB.DeleteSession(sessionID)
		if err != nil {
			log.Warnf("failed to delete session: %s", err.Error())
		}
		log.Debug("session deleted: ", sessionID)
	})
}

func CleanupSessions() {
	expiryTimes, err := dal.DB.GetSessionExpiries()
	if err != nil {
		log.WithError(err).Error("failed to get expiry times")
	}
	log.Infof("scheduling the deletion of %d sessions", len(expiryTimes))

	for k, v := range expiryTimes {
		var ttl int64
		if ttl = int64(v) - time.Now().Unix(); ttl < 0 {
			ttl = 0
		}
		endSessionAfterDuration(k, time.Second*time.Duration(ttl))
	}
}
