package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/modules/repository"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type CredentialService struct {
	db         *gorm.DB
	repository repository.Repository
}

func NewCredentialService(db *gorm.DB, repo repository.Repository) *CredentialService {
	return &CredentialService{
		db:         db,
		repository: repo,
	}
}

const cost = 10

func (c *CredentialService) CreateCredential(credential *credential.Credentials) error {

	if len(strings.TrimSpace(credential.Email)) == 0 || len(strings.TrimSpace(credential.Password)) == 0 {
		return apperror.NewValidationError("EMPTY_VALUE", "email or password cannot be empty")
	}

	hashedPassword, err := hashPassword(credential.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	credential.Password = string(hashedPassword)

	if err := c.db.Create(credential).Error; err != nil {
		return apperror.NewDatabaseError("Error creating credentials")
	}

	// err = CredentialService.Add(uow, newUser)
	// if err != nil {
	// 	return apperror.NewDatabaseError("Failed to create admin user")
	// }

	return nil
}

func hashPassword(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func CheckPassword(userPassword string, inputPassword string) error {

	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword))

}
