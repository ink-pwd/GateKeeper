package security

import "golang.org/x/crypto/bcrypt"

func CryptoHash(password string) (string, error) {
	var (
		password_hash []byte
		err           error
	)
	password_hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(password_hash), nil
}

func Compare(password, password_hash string) bool {
	var (
		err error
	)
	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
