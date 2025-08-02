package service

import (
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/models/user"
	"errors"

	"github.com/jinzhu/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(DB *gorm.DB) *UserService {
	return &UserService{
		db: DB,
	}
}

func (u *UserService) CreateAdmin(fname, lname, email, password string) (*user.User, error) {

	foundUser, err := credential.CreateCredential(email, password)
	if err != nil {
		return nil, err
	}
	user := user.CreateAdmin(fname, lname, foundUser)

	err = u.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) Login(email, password string) (*user.User, error) {

	foundUser, foundCredentials, err := u.FindCredential(email)
	if err != nil {
		return nil, err
	}

	if foundCredentials == nil {
		return nil, errors.New("user credentials not found")
	}

	err = credential.CheckPassword(foundCredentials.Password, password)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (u *UserService) FindCredential(email string) (*user.User, *credential.Credentials, error) {

	foundCredential := credential.Credentials{}
	err := u.db.Where("email = ?", email).First(&foundCredential).Error
	if err != nil {
		return nil, nil, err
	}

	foundUser := user.User{}
	err = u.db.Preload("Credentials").Preload("Contacts").Where("credentials_id = ?", foundCredential.ID).First(&foundUser).Error
	if err != nil {
		return nil, nil, err
	}

	if foundUser.ID == 0 {
		return nil, nil, errors.New("user not found")
	}

	return &foundUser, &foundCredential, nil
}
