package reldb

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	VFirstName *string `db:"vFirstName" json:"vFirstName,omitempty"`
	VLastName  *string `db:"vLastName" json:"vLastName,omitempty"`
	VEmail     *string `db:"vEmail" json:"vEmail,omitempty"`
	VRole      *string `db:"vRole" json:"vRole,omitempty"`
}
type User struct {
	IUserID uint `db:"iUserID" json:"-"`
	UserData
	VPasswordHash *string `db:"vPasswordHash" json:"vPasswordHash,omitempty"`
}

type LoggedInUser struct {
	UserData
	JWTToken     *string `json:"jwtToken,omitempty"`
	SecureCookie *string `json:"secureCookie,omitempty"`
}

func generateHash(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error generating hashed/salted password: %s", err)
	}

	return string(hash), nil
}

func verifyPassword(storedPassword, sentPassword string) bool {

	bStoredPassword := []byte(storedPassword)
	bSentPassword := []byte(sentPassword)

	if err := bcrypt.CompareHashAndPassword(bStoredPassword, bSentPassword); err != nil {
		return false
	}

	return true
}

func (m *Model) UserIsAuthorized(email string) (bool, error) {

	query := `SELECT COUNT(*) FROM admin_users WHERE vEmail = ?`
	var count int
	if err := m.Get(&count, query, email); err != nil {
		return false, err
	}

	if count != 1 {
		return false, fmt.Errorf("email %s does not exist (count: %d)", email, count)
	}

	return true, nil
}

func (m *Model) UpdateUser(user User) (int64, error) {

	if user.VPasswordHash == nil || *user.VPasswordHash == "" ||
		user.VEmail == nil || *user.VEmail == "" {
		return 0,
			fmt.Errorf("email and password are required and cannot be empty")
	}

	pwHash, err := generateHash(*user.VPasswordHash)
	if err != nil {
		return 0, fmt.Errorf("error generating Hashed password: %s", err)
	}

	query := `UPDATE admin_users SET
				vFirstName = ?,
				vLastName = ?,
				vPasswordHash = ?
			WHERE vEmail = ?`

	result, err := m.Exec(query,
		*user.VFirstName,
		*user.VLastName,
		pwHash,
		*user.VEmail,
	)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *Model) Login(email, password string) (User, error) {

	var user User
	query := `SELECT
				iUserID,
				vFirstName,
				vLastName,
				vEmail,
				vRole,
				vPasswordHash
			FROM admin_users 
			WHERE vEmail = ?`

	if err := m.Get(&user, query, email); err != nil {
		return User{}, err
	}

	if !verifyPassword(*user.VPasswordHash, password) {
		return User{}, fmt.Errorf("incorrect password")
	}

	return user, nil
}

func (m *Model) SaveCookie(userID uint, sessionID, cookieStr string, expires time.Time) error {

	delQuery := `DELETE FROM session WHERE iUserID = ?`
	insQuery := `INSERT INTO session 
				(vSessionID, iUserID, vCookieStr, dtExpires) 
			VALUES (?, ?, ?, ?)`

	tx, err := m.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(delQuery, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(insQuery, sessionID, userID, cookieStr, expires)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (m *Model) RetrieveCookieString(id string, iUserID uint) (string, error) {

	var cookieStr string
	query := `SELECT vCookieStr FROM session WHERE vSessionID = ? AND iUserID = ?`
	if err := m.Get(&cookieStr, query, id, iUserID); err != nil {
		return "", err
	}

	return cookieStr, nil
}

func (m *Model) RetrieveUserByID(iUserID uint) (UserData, error) {

	var userData UserData
	query := `SELECT
				vFirstName,
				vLastName,
				vEmail,
				vRole
			FROM admin_users 
			WHERE iUserID = ?`

	if err := m.Get(&userData, query, iUserID); err != nil {
		return userData, err
	}

	return userData, nil
}

func (m *Model) VerifyActiveUser(vEmail, vRole string) error {

	query := `SELECT COUNT(*) FROM admin_users a
				JOIN session s ON a.iUserID = s.iUserID
				WHERE 
					a.vEmail = ? AND a.vRole = ?`
	var count int
	if err := m.Get(&count, query, vEmail, vRole); err != nil {
		return err
	}

	if count != 1 {
		return fmt.Errorf("user has %d active sessions", count)
	}

	return nil
}
