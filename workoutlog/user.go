package workoutlog

import (
	"fmt"

	"github.com/scottbrodersen/homegym/auth"

	"github.com/scottbrodersen/homegym/dal"
)

type User struct {
	ID             string    `json:"id"`
	Role           auth.Role `validate:"required" json:"role"`
	Email          string    `validate:"required" json:"email"`
	PwdHash        string    `json:"-"`
	PwdHashVersion string    `json:"-"`
}

var FrontDesk UserAdmin = new(userManager)

type UserAdmin interface {
	NewUser(username string, role auth.Role, email, password string) (*User, error)
	GetUser(username string) (*User, error)
}

type userManager struct{}

func (u *userManager) NewUser(username string, role auth.Role, email, password string) (*User, error) {
	pwdUtil := auth.PasswordUtil{Version: auth.LatestVersion()}
	hash, err := pwdUtil.Hash(password)
	if err != nil {
		return nil, err
	}

	if err = dal.DB.NewUser(username, email, hash, auth.LatestVersion(), string(role)); err != nil {
		return nil, err
	}

	user := User{
		ID:             username,
		Role:           role,
		Email:          email,
		PwdHash:        hash,
		PwdHashVersion: auth.LatestVersion(),
	}

	return &user, nil
}

func (u *userManager) GetUser(username string) (*User, error) {
	email, pwdHash, pwdHashVersion, role, err := dal.DB.ReadUser(username)
	if err != nil {
		return nil, fmt.Errorf("failed to read user: %w", err)
	}

	if email == nil && pwdHash == nil && pwdHashVersion == nil && role == nil {
		return nil, fmt.Errorf("user not found")
	}

	var user User = User{
		ID:             username,
		Email:          *email,
		PwdHash:        *pwdHash,
		PwdHashVersion: *pwdHashVersion,
		Role:           auth.Role(*role),
	}

	return &user, nil
}

// Updates tables with field values of u
// Does not update the role, password, or password version
func (u *User) UpdateUserProfile() error {
	_, err := FrontDesk.GetUser(u.ID)
	if err != nil {
		return err
	}

	if err := dal.DB.UpdateUserProfile(u.ID, u.Email); err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateUserPassword(newPassword string) error {
	user, err := FrontDesk.GetUser(u.ID)
	if err != nil {
		return err
	}

	// check that the new password is different
	pwdUtil := auth.PasswordUtil{Version: user.PwdHashVersion}
	newPwdHash, err := pwdUtil.Hash(newPassword)
	if err != nil {
		return err
	}

	if newPwdHash == user.PwdHash {
		return fmt.Errorf("password unchanged")
	}

	// Use the latest password version when password is updated
	pwdUtil.Version = auth.LatestVersion()
	newHash, err := pwdUtil.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := dal.DB.UpdateUserPassword(u.ID, newHash, auth.LatestVersion()); err != nil {
		return err
	}

	u.PwdHash = newHash
	u.PwdHashVersion = auth.LatestVersion()

	return nil
}

func (u *userManager) ChangeUserRole(userID string, role auth.Role) error {
	user, err := FrontDesk.GetUser(userID)
	if err != nil {
		return err
	}

	if user.Role == role {
		return fmt.Errorf("role unchanged")
	}

	if err := dal.DB.ChangeUserRole(userID, string(role)); err != nil {
		return err
	}

	return nil
}
