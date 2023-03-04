package sphelper

import (
	"os"
	"testing"
)

func TestCredentials(t *testing.T) {
	// store current credentials
	username := os.Getenv("SP_USERNAME")
	password := os.Getenv("SP_PASSWORD")

	// arrange
	os.Setenv("SP_USERNAME", "user")
	os.Setenv("SP_PASSWORD", "password")

	// act
	cred, err := GetCredentials()

	// assert
	if err != nil {
		t.Errorf("could not get credentials: %v", err)
	}
	expected := Credentials{
		Username:    "user",
		Password:    "password",
		PasswordMd5: "5f4dcc3b5aa765d61d8327deb882cf99",
	}
	if cred != expected {
		t.Errorf("%v != %v", cred, expected)
	}

	// restore credentials
	os.Setenv("SP_USERNAME", username)
	os.Setenv("SP_PASSWORD", password)
}

func TestMissingEnv(t *testing.T) {
	// store current credentials
	username := os.Getenv("SP_USERNAME")
	password := os.Getenv("SP_PASSWORD")

	// arrange
	os.Unsetenv("SP_USERNAME")
	os.Unsetenv("SP_PASSWORD")

	// act
	cred, err := GetCredentials()

	// assert
	if err != userNameOrPasswordNotFound {
		t.Errorf("%v != %v", err, userNameOrPasswordNotFound)
	}
	expected := Credentials{}
	if cred != expected {
		t.Errorf("%v != %v", cred, expected)
	}

	// restore credentials
	os.Setenv("SP_USERNAME", username)
	os.Setenv("SP_PASSWORD", password)
}
