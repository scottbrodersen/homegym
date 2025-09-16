package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scottbrodersen/homegym/dal"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// A Session represents the state of a user session.
type Session struct {
	UserID    string
	SessionID string
	Expiry    int64
}

type Role string

// type CryptoKey struct {
// 	Key     *rsa.PrivateKey
// 	Retired bool
// }

var Admin Role = "admin"
var User Role = "user"

// keyRotationTime is the default time of daily key rotation in format "hh:mm:ss".
var keyRotationTime = "10:00:00"    // default, UTC
var secondsToDelayKeyDeletion = 600 // 10 mins
const tokenLifeInSeconds = 1800     // 30 mins
const sessionLifeInMinutes = 1440   // 24 hours

// SetKeyRotationTime schedules the time of the daily key rotation.
// The time argument is expressed as hh:mm:ss
func SetKeyRotationTime(rotationTime string) error {
	rx := regexp.MustCompile("^[0-9]{2}:[0-9]{2}:[0-9]{2}$")
	if !rx.Match([]byte(rotationTime)) {
		err := fmt.Errorf("invalid time, using default: %s", keyRotationTime)
		slog.Error("key rotation time not in hh:mm:ss format", "error", err.Error())
		slog.Info("using default rotation time.")
		return err
	}
	slog.Info(fmt.Sprintf("key rotation time set to: %s", rotationTime))
	keyRotationTime = rotationTime
	return nil
}

// A KeyUsage defines the different uses of keys.
type KeyUsage struct {
	Token  string
	Pepper string
}

var KeyTypes KeyUsage = KeyUsage{
	Token:  "token",
	Pepper: "pepper",
}

// InitiateKeyRotation starts the daily key rotation process.
// Rotates keys immediately if no usable key is found.
func InitiateKeyRotation(keyType string) error {
	// get current key and if there is none create one now
	if keyType == KeyTypes.Token {
		signingKey, err := getSigningKey()
		if err != nil {
			slog.Error("failed to get signing key -- ignoring", "error", err.Error())
		}
		if signingKey == nil {
			slog.Warn("signing key not found; initiating key rotation immediately")
			rotateKey(KeyTypes.Token)
		}
	}

	slog.Info("initiating daily key rotation.")
	dailyKeyRotation(keyType)

	return nil
}

// getSigningKey returns the private key used for signing tokens.
// Gets all token-signing keys and returns the (first) active key.
// Returns nil when no active key is found.
func getSigningKey() (*rsa.PrivateKey, error) {
	activeKeys, _, err := dal.DB.GetKeys(KeyTypes.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to read keys: %w", err)
	}

	for _, key := range activeKeys {

		keyRSA, err := byteSliceToRSAKey(key)
		if err != nil {
			return nil, err
		}

		return keyRSA, nil
	}
	return nil, nil
}

// dailyKeyRotation schedules the next key rotation.
// Calculates the time until the next rotation and schedules the rotation execution.
// This function is recursive so that once a rotation occurs, the next rotation is scheduled.
func dailyKeyRotation(keyType string) {
	// fallback
	rotationDelay := time.Second * 0

	now := time.Now().UTC()

	h, m, s := now.Clock()
	seconds := h*3600 + m*60 + s

	slog.Debug(fmt.Sprintf("Key rotation time is set to %s", keyRotationTime))

	hmsSched := strings.Split(keyRotationTime, ":")
	hSched, errH := strconv.Atoi(hmsSched[0])
	mSched, errM := strconv.Atoi(hmsSched[1])
	sSched, errS := strconv.Atoi(hmsSched[2])

	if errH != nil || errM != nil || errS != nil {
		slog.Warn("error parsing key rotation time -- using fallback delay")
	} else {
		secondsSched := hSched*3600 + mSched*60 + sSched
		secondsToRot := secondsSched - seconds

		// if the scheduled time is in the past, do it tomorrow
		if secondsToRot <= 0 {
			secondsToRot += 24 * 60 * 60
		}
		rotationDelay = time.Second * time.Duration(secondsToRot)
	}

	slog.Info(fmt.Sprintf("Scheduling %s key rotation to occur in %s", keyType, rotationDelay.String()))

	time.AfterFunc(rotationDelay, func() {
		err := rotateKey(keyType)
		if err != nil {
			slog.Error("failed to rotate key", "error", err.Error())
		}
		dailyKeyRotation(keyType)
	})
}

// rotateKey rotates the key of a specific type.
// Old keys are scheduled for deletion.
func rotateKey(keyType string) error {
	slog.Info(fmt.Sprintf("rotating %s key", keyType))

	keyID := uuid.NewString()
	newKey, err := generateRsaPrivate(2048)
	if err != nil {
		return fmt.Errorf("failed to create key: %w", err)
	}

	newKeyByte := rsaKeyToByteSlice(newKey)

	err = dal.DB.RotateKeys(newKeyByte, keyID, keyType)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("%s key created; id %s", keyType, keyID))

	// delete old key(s)
	_, retiredKeys, err := dal.DB.GetKeys(keyType)
	if err != nil {
		return fmt.Errorf("failed to get keys during rotation: %w", err)
	}

	for id := range retiredKeys {
		deleteKeyAfterDuration(id, keyType, time.Second*time.Duration(secondsToDelayKeyDeletion))
	}

	return nil
}

func deleteKeyAfterDuration(keyID, usage string, duration time.Duration) {
	slog.Info(fmt.Sprintf("deletion of %s key scheduled in %s; id: %s", usage, duration.String(), keyID))

	time.AfterFunc(duration, func() {

		err := dal.DB.DeleteKey(keyID, usage)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to delete %s key", usage), "error", err.Error())
		}
		slog.Info(fmt.Sprintf("deleted %s key on schedule; id: %s", usage, keyID))

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

// A hashPasswordVersion returns the name of a password version.
type hashPasswordVersion func(string) (string, error)

// A validatePasswordVersion validates a password.
// An error is returned when validation fails.
type validatePasswordVersion func(string, string) error

// A PasswordUtil hashes and validates passwords.
// Version determines the version of the hashing and validation functions used.
type PasswordUtil struct {
	Version string
}

var (
	errPwdShort = "password too short"
	errPwdLong  = "password too long"
	errPwdWrong = "bad credentials"
)

// hashFunctions maps password version names to a hashing functions.
// Add new versions as needed.
// keep old function versions to use with non-updated credential versions
var hashFunctions map[string]hashPasswordVersion = map[string]hashPasswordVersion{
	"v1":   hashPasswordVersion1,
	"mock": hashPasswordMock,
}

// validateFunctions maps password version names to password validation functions.
var validateFunctions map[string]validatePasswordVersion = map[string]validatePasswordVersion{
	"v1":   validatePasswordVersion1,
	"mock": validatePasswordMock,
}

// LatestVersion returns the current version of hashing and validation functions that we're using.
// This version should be used with new credentials.
func LatestVersion() string {
	return "v1"
}

// Hash returns the hash of a password.
// Uses the hashing version that the PasswordUtil contains
func (p *PasswordUtil) Hash(pwd string) (string, error) {
	if hashFunc, ok := hashFunctions[p.Version]; ok {
		return hashFunc(pwd)
	}
	return "", fmt.Errorf("invalid version")
}

// ValidatePassword validates a password.
// Uses the validation version that the PasswordUtil contains.
func (p *PasswordUtil) ValidatePassword(pwd string, hash string) error {
	if validateFunc, ok := validateFunctions[p.Version]; ok {
		return validateFunc(pwd, hash)
	}
	return fmt.Errorf("invalid version")
}

// The password hashing function for v1.
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

// The password validation function for v1.
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

// An Authorizer issues and validates tokens.
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

// A Claims validates token claims.
type Claims struct {
	jwt.RegisteredClaims
	GymClaims
}

// A GymClaims contains custom token claims.
type GymClaims struct {
	Role string
}

// Valid validates the claims in a token.
func (c *Claims) Valid() error {
	if c.Issuer != issuer {
		slog.Debug(fmt.Sprintf("token claim contained invalid issuer: %s", c.Issuer))
		return fmt.Errorf("invalid issuer")
	}
	return nil
}

func tokenExpiryTime() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(tokenLifeInSeconds)))
}

// IssueToken authenticates user credentials.
// Creates a token.
// Stores the session.
// Returns the token string and sessionID.
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
	slog.Info("Session created ", "id", sessionID)
	endSessionAfterDuration(sessionID, time.Hour*time.Duration(24))

	return &tokenStr, &sessionID, nil
}

// ValidateToken  checks whether a session is valid.
// After validating the user's token, returns a new token with updated expiry.
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

// TokenClaims returns the claims in a JWT.
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
	slog.Info("Scheduling deletion of session", "in", duration)
	time.AfterFunc(duration, func() {
		err := dal.DB.DeleteSession(sessionID)
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to delete session: %s", err.Error()))
		}
		slog.Debug("session deleted. ", "sessionID", sessionID)
	})
}

// CleanupSessions removes expired sessions.
func CleanupSessions() {
	expiryTimes, err := dal.DB.GetSessionExpiries()
	if err != nil {
		slog.Error("failed to get expiry times", "error", err.Error())
	}
	slog.Info(fmt.Sprintf("scheduling the deletion of %d sessions", len(expiryTimes)))

	for k, v := range expiryTimes {
		var ttl int64
		if ttl = int64(v) - time.Now().Unix(); ttl < 0 {
			ttl = 0
		}
		endSessionAfterDuration(k, time.Second*time.Duration(ttl))
	}
}
