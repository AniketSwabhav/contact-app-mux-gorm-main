package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/models/user"
	"strings"

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

	if strings.TrimSpace(fname) == "" {
		return nil, apperror.NewMissingFieldsError("first name is required")
	}

	if strings.TrimSpace(lname) == "" {
		return nil, apperror.NewMissingFieldsError("last name is required")
	}

	if strings.TrimSpace(email) == "" {
		return nil, apperror.NewMissingFieldsError("email is required")
	}

	if strings.TrimSpace(password) == "" {
		return nil, apperror.NewMissingFieldsError("password is required")
	}

	var existing credential.Credentials
	if err := u.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, apperror.NewValidationError("DUPLICATE_EMAIL", "Email already in use")
	}

	foundUser, err := credential.CreateCredential(email, password)
	if err != nil {
		return nil, apperror.NewValidationError("INVALID_CREDENTIAL", err.Error())
	}

	user := user.CreateAdmin(fname, lname, foundUser)

	err = u.db.Create(user).Error
	if err != nil {
		return nil, apperror.NewDatabaseError("Failed to create admin user")
	}

	return user, nil
}

func (u *UserService) CreateUser(fname, lname, email, password string) (*user.User, error) {

	if strings.TrimSpace(fname) == "" {
		return nil, apperror.NewMissingFieldsError("first name is required")
	}

	if strings.TrimSpace(lname) == "" {
		return nil, apperror.NewMissingFieldsError("last name is required")
	}

	if strings.TrimSpace(email) == "" {
		return nil, apperror.NewMissingFieldsError("email is required")
	}

	if strings.TrimSpace(password) == "" {
		return nil, apperror.NewMissingFieldsError("password is required")
	}

	var existing credential.Credentials
	if err := u.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, apperror.NewValidationError("DUPLICATE_EMAIL", "Email already in use")
	}

	foundUser, err := credential.CreateCredential(email, password)
	if err != nil {
		return nil, apperror.NewValidationError("INVALID_CREDENTIAL", err.Error())
	}

	user := user.CreateUser(fname, lname, foundUser)

	err = u.db.Create(user).Error
	if err != nil {
		return nil, apperror.NewDatabaseError("Failed to create admin user")
	}

	return user, nil
}

func (u *UserService) Login(email, password string) (*user.User, error) {

	foundUser, foundCredentials, err := u.FindCredential(email)
	if err != nil {
		return nil, apperror.NewNotFoundError(err.Error())
	}

	if foundCredentials == nil {
		return nil, apperror.NewNotFoundError("user credentials not found")
	}

	err = credential.CheckPassword(foundCredentials.Password, password)
	if err != nil {
		return nil, apperror.NewInValidPasswordError("password is incorrect")
	}

	return foundUser, nil
}

func (u *UserService) FindCredential(email string) (*user.User, *credential.Credentials, error) {

	foundCredential := credential.Credentials{}
	err := u.db.Where("email = ?", email).First(&foundCredential).Error
	if err != nil {
		return nil, nil, apperror.NewNotFoundError("user not found with given email")
	}

	foundUser := user.User{}
	err = u.db.Preload("Credentials").Preload("Contacts").Where("credentials_id = ?", foundCredential.ID).First(&foundUser).Error
	if err != nil {
		return nil, nil, apperror.NewNotFoundError("Invalid credentials_id")
	}

	if foundUser.ID == 0 {
		return nil, nil, apperror.NewNotFoundError("user does not exists")
	}

	return &foundUser, &foundCredential, nil
}

func (u *UserService) GetAllUsers() (*[]user.User, error) {

	var allUsers = &[]user.User{}

	err := u.db.Preload("Credentials").Preload("Contacts").Find(&allUsers).Error
	if err != nil {
		return nil, apperror.NewDatabaseError("unable to fetch users")
	}

	return allUsers, nil
}

func (u *UserService) GetUsersPaginated(page, pageSize int) (*[]user.User, error) {
	var users = &[]user.User{}

	offset := (page - 1) * pageSize

	err := u.db.Preload("Credentials").Preload("Contacts").Limit(pageSize).Offset(offset).Find(users).Error

	if err != nil {
		return nil, apperror.NewDatabaseError("Unable to fetch users")
	}

	return users, nil
}

func (u *UserService) Get(userID string) (*user.User, error) {

	var foundUser = &user.User{}

	err := u.db.Preload("Credentials").Preload("Contacts").Where("user_id = ?", userID).First(&foundUser).Error
	if err != nil {
		return nil, apperror.NewNotFoundError("User with given id not found")
	}

	return foundUser, nil
}

func (u *UserService) Update(userId, firstName, lastName, email string) (*user.User, error) {

	var existingUser user.User

	err := u.db.Preload("Credentials").Where("user_id = ?", userId).First(&existingUser).Error
	if err != nil {
		return nil, apperror.NewNotFoundError("user not found")
	}

	existingUser.FirstName = firstName
	existingUser.LastName = lastName

	if existingUser.Credentials != nil && existingUser.Credentials.Email != email {
		var count int64
		u.db.Model(&credential.Credentials{}).Where("email = ?", email).Count(&count)
		if count > 0 {
			return nil, apperror.NewDuplicateEntryError("email already in use")
		}
		existingUser.Credentials.Email = email
	}

	err = u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existingUser).Error; err != nil {
			return err
		}
		if existingUser.Credentials != nil {
			if err := tx.Save(&existingUser.Credentials).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, apperror.NewDatabaseError("Failed to update user")
	}

	return &existingUser, nil
}

func (u *UserService) Delete(userId string) error {
	var existingUser user.User

	err := u.db.Where("user_id = ?", userId).First(&existingUser).Error
	if err != nil {
		return apperror.NewNotFoundError("user not found")
	}

	if err := u.db.Delete(&existingUser).Error; err != nil {
		return apperror.NewDatabaseError("Failed to delete user")
	}

	return nil
}
