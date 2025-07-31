package credential

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// to store a map of credentials
// var credentialStore = make(map[string]*Credentials)

const cost = 10

type Credentials struct {
	Email    string `json:"Email" gorm:"unique;not null;type:varchar(100)"`
	Password string `json:"Password" gorm:"not null;type:varchar(100)"`
}

func CreateCredential(email string, password string) (*Credentials, error) {

	if len(strings.TrimSpace(email)) == 0 || len(strings.TrimSpace(password)) == 0 {
		return nil, errors.New("credentials cannot be empty")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	// credentialId := uuid.New()

	newCredential := &Credentials{
		// CredentialID: credentialId.String(),
		Email:    email,
		Password: string(hashedPassword),
	}

	// credentialStore[credentialId.String()] = newCredential

	return newCredential, nil
}

func hashPassword(password string) ([]byte, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
