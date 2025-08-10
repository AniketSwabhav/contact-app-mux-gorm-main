package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/models/user"
	"contact_app_mux_gorm_main/modules/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const cost = 10

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

	if newUser.Credentials == nil {
		return apperror.NewValidationError("INVALID_INPUT", "Missing credentials")
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	newUser.IsAdmin = true

	hashedPassword, err := hashPassword(newUser.Credentials.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser.Credentials.Password = string(hashedPassword)

	err = uow.DB.Create(newUser).Error
	if err != nil {
		return apperror.NewDatabaseError("Failed to create user")
	}

	uow.Commit()
	return nil
}

func (service *UserService) CreateUser(newUser *user.User) error {

	err := service.doesEmailExists(newUser.Credentials.Email)
	if err != nil {
		return err
	}

	if newUser.Credentials == nil {
		return apperror.NewValidationError("INVALID_INPUT", "Missing credentials")
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	hashedPassword, err := hashPassword(newUser.Credentials.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser.Credentials.Password = string(hashedPassword)

	err = uow.DB.Create(newUser).Error
	if err != nil {
		return apperror.NewDatabaseError("Failed to create user")
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
	err = uow.DB.Preload("Credentials").
		Where("id = ?", foundCredential.UserID).First(&foundUser).Error

	if err != nil {
		return apperror.NewDatabaseError("Could not retrieve user")
	}

	*claim = authorization.Claims{
		UserID:   foundUser.ID,
		IsAdmin:  foundUser.IsAdmin,
		IsActive: foundUser.IsActive,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Minute).Unix(),
		},
	}
	uow.Commit()
	return nil
}

func (service *UserService) GetAllUsers(allUsers *[]user.UserDTO, totalCount *int) error {

	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()

	limit := 5
	offset := 0

	err := service.repository.GetAll(uow, allUsers, repository.PreloadAssociations([]string{"Credentials", "Contacts"}), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		return err
	}

	err = service.repository.GetCount(uow, allUsers, totalCount)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

func (service *UserService) GetUserByID(targetUser *user.UserDTO) error {

	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()

	err := service.repository.GetRecordByID(uow, targetUser.ID, targetUser, repository.PreloadAssociations([]string{"Credentials", "Contacts"}))
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) UpdateUser(userToUpdate *user.User) error {
	err := service.doesUserExist(userToUpdate.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	tempUser := user.User{}

	err = service.repository.GetRecordByID(uow, userToUpdate.ID, &tempUser, repository.Select("`created_at`"),
		repository.Filter("`id` = ?", userToUpdate.ID))
	if err != nil {
		return err
	}

	err = service.repository.UpdateWithMap(uow, userToUpdate, map[string]interface{}{
		"first_name": userToUpdate.FirstName,
		"last_name":  userToUpdate.LastName,
		"is_admin":   userToUpdate.IsAdmin,
		"is_active":  userToUpdate.IsActive,
	})

	uow.Commit()
	return nil
}

func (service *UserService) Delete(userToDelete *user.User) error {

	err := service.doesUserExist(userToDelete.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	if err := service.repository.UpdateWithMap(uow, userToDelete, map[string]interface{}{
		"DeletedAt": time.Now(),
	},
		repository.Filter("`id`=?", userToDelete.ID)); err != nil {
		uow.RollBack()
		return err
	}

	if err := uow.DB.Where("user_id = ?", userToDelete.ID).Delete(&credential.Credentials{}).Error; err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

func (service *UserService) doesEmailExists(Email string) error {
	exists, _ := repository.DoesEmailExist(service.db, Email, credential.Credentials{},
		repository.Filter("`email` = ?", Email))
	if exists {
		return apperror.NewValidationError("EMAIL_ALREADY_EXISTS", "Email is already registered")
	}
	return nil
}

func (service *UserService) doesUserExist(ID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForUser(service.db, ID, user.User{},
		repository.Filter("`id` = ?", ID))
	if !exists || err != nil {
		return apperror.NewValidationError("INVALID_ID", "User ID is Invalid")
	}
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

	return &foundUser, &foundCredential, nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}
