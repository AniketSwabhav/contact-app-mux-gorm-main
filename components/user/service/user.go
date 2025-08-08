package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	credentialServices "contact_app_mux_gorm_main/components/credential/service"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/models/user"
	"contact_app_mux_gorm_main/modules/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db         *gorm.DB
	repository repository.Repository
}

func NewUserService(DB *gorm.DB, repo repository.Repository) *UserService {
	return &UserService{
		db:         DB,
		repository: repo,
	}
}

func (service *UserService) CreateAdmin(newUser *user.User) error {

	err := service.doesEmailExists(newUser.Credentials.Email)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	newUser.UserID = uuid.New().String()
	newUser.IsAdmin = true

	credential := credential.Credentials{
		Email:    newUser.Credentials.Email,
		Password: newUser.Credentials.Password,
	}
	newCredentials := credentialServices.NewCredentialService(uow.DB, service.repository)

	err = newCredentials.CreateCredential(&credential)
	if err != nil {
		return apperror.NewDatabaseError(err.Error())
	}

	newUser.Credentials = &credential

	err = service.repository.Add(uow, newUser)
	if err != nil {
		return apperror.NewDatabaseError("Failed to create admin user")
	}

	uow.Commit()
	return nil
}

func (service *UserService) CreateUser(newUser *user.User) error {

	err := service.doesEmailExists(newUser.Credentials.Email)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	newUser.UserID = uuid.New().String()

	credential := credential.Credentials{
		Email:    newUser.Credentials.Email,
		Password: newUser.Credentials.Password,
	}
	newCredentials := credentialServices.NewCredentialService(uow.DB, service.repository)

	newCredentials.CreateCredential(&credential)

	newUser.Credentials = &credential

	err = service.repository.Add(uow, newUser)
	if err != nil {
		return apperror.NewDatabaseError("Failed to create admin user")
	}

	uow.Commit()
	return nil
}

func (service *UserService) Login(userCredential *credential.Credentials, claim *authorization.Claims) error {

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	exists, err := repository.DoesEmailExist(service.db, userCredential.Email, credential.Credentials{},
		repository.Filter("`email` = ?", userCredential.Email))
	if err != nil {
		return apperror.NewDatabaseError("Error checking if email exists")
	}
	if !exists {
		return apperror.NewNotFoundError("Email not found")
	}

	foundCredential := credential.Credentials{}
	err = uow.DB.Where("email = ?", userCredential.Email).First(&foundCredential).Error
	if err != nil {
		return apperror.NewDatabaseError("Could not retrieve credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundCredential.Password), []byte(userCredential.Password))
	if err != nil {
		return apperror.NewInValidPasswordError("Incorrect password")
	}

	foundUser := user.User{}
	err = uow.DB.Preload("Credentials").Preload("Contacts").
		Where("credentials_id = ?", foundCredential.ID).First(&foundUser).Error

	if err != nil {
		return apperror.NewDatabaseError("Could not retrieve user")
	}

	*claim = authorization.Claims{
		UserID:   foundUser.UserID,
		IsAdmin:  foundUser.IsAdmin,
		IsActive: foundUser.IsActive,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Minute).Unix(),
		},
	}
	uow.Commit()
	return nil
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

func (service *UserService) doesEmailExists(Email string) error {
	exists, err := repository.DoesEmailExist(service.db, Email, credential.Credentials{},
		repository.Filter("`email` = ?", Email))
	if err != nil {
		return apperror.NewDatabaseError("Error checking email existence")
	}
	if exists {
		return apperror.NewValidationError("EMAIL_ALREADY_EXISTS", "Email is already registered")
	}
	return nil
}

// func (service *UserService) doesUserExist(ID uint) error {
// 	exists, err := repository.DoesRecordExistForUser(service.db, ID, user.User{},
// 		repository.Filter("`id` = ?", ID))
// 	if !exists || err != nil {
// 		return errors.NewValidationError("User ID is Invalid")
// 	}
// 	return nil
// }

// func (u *UserService) GetAllUsers() (*[]user.User, error) {

// 	var allUsers = &[]user.User{}

// 	err := u.db.Preload("Credentials").Preload("Contacts").Find(&allUsers).Error
// 	if err != nil {
// 		return nil, apperror.NewDatabaseError("unable to fetch users")
// 	}

// 	return allUsers, nil
// }

// func (u *UserService) GetUsersPaginated(page, pageSize int) (*[]user.User, error) {
// 	var users = &[]user.User{}

// 	offset := (page - 1) * pageSize

// 	err := u.db.Preload("Credentials").Preload("Contacts").Limit(pageSize).Offset(offset).Find(users).Error

// 	if err != nil {
// 		return nil, apperror.NewDatabaseError("Unable to fetch users")
// 	}

// 	return users, nil
// }

// func (u *UserService) Get(userID string) (*user.User, error) {

// 	var foundUser = &user.User{}

// 	err := u.db.Preload("Credentials").Preload("Contacts").Where("user_id = ?", userID).First(&foundUser).Error
// 	if err != nil {
// 		return nil, apperror.NewNotFoundError("User with given id not found")
// 	}

// 	return foundUser, nil
// }

// func (u *UserService) Update(userId, firstName, lastName, email string) (*user.User, error) {

// 	var existingUser user.User

// 	err := u.db.Preload("Credentials").Where("user_id = ?", userId).First(&existingUser).Error
// 	if err != nil {
// 		return nil, apperror.NewNotFoundError("user not found")
// 	}

// 	existingUser.FirstName = firstName
// 	existingUser.LastName = lastName

// 	if existingUser.Credentials != nil && existingUser.Credentials.Email != email {
// 		var count int64
// 		u.db.Model(&credential.Credentials{}).Where("email = ?", email).Count(&count)
// 		if count > 0 {
// 			return nil, apperror.NewDuplicateEntryError("email already in use")
// 		}
// 		existingUser.Credentials.Email = email
// 	}

// 	err = u.db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Save(&existingUser).Error; err != nil {
// 			return err
// 		}
// 		if existingUser.Credentials != nil {
// 			if err := tx.Save(&existingUser.Credentials).Error; err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, apperror.NewDatabaseError("Failed to update user")
// 	}

// 	return &existingUser, nil
// }

// func (u *UserService) Delete(userId string) error {
// 	var existingUser user.User

// 	err := u.db.Where("user_id = ?", userId).First(&existingUser).Error
// 	if err != nil {
// 		return apperror.NewNotFoundError("user not found")
// 	}

// 	if err := u.db.Delete(&existingUser).Error; err != nil {
// 		return apperror.NewDatabaseError("Failed to delete user")
// 	}

// 	return nil
// }
