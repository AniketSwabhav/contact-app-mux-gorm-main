package service

import (
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/models/user"

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

func (s *UserService) CreateAdmin(fname, lname, email, password string) (*user.User, error) {

	newCredential, err := credential.CreateCredential(email, password)
	if err != nil {
		return nil, err
	}
	user := user.CreateAdmin(fname, lname, newCredential)

	err = s.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
