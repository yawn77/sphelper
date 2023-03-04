package sphelper

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// errors
const (
	userNameOrPasswordNotFound = Error("environment variable SP_USERNAME or SP_PASSWORD not set")
)

type Credentials struct {
	Username    string
	Password    string
	PasswordMd5 string
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetCredentials() (Credentials, error) {
	username := os.Getenv("SP_USERNAME")
	password := os.Getenv("SP_PASSWORD")
	if username == "" || password == "" {
		return Credentials{}, userNameOrPasswordNotFound
	}
	return Credentials{
		Username:    username,
		Password:    password,
		PasswordMd5: GetMD5Hash(password),
	}, nil
}
