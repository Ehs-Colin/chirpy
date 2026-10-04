package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hash_password, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash_password, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return result, nil
}
