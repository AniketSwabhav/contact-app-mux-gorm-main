package service

import "github.com/jinzhu/gorm"

type CredentialService struct {
	db *gorm.DB
}

func NewCredentialService(db *gorm.DB) *CredentialService {
	return &CredentialService{
		db: db,
	}
}
