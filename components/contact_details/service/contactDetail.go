package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/models/contactdetail"
	"contact_app_mux_gorm_main/modules/repository"

	"github.com/jinzhu/gorm"
)

type ContactDeatailsService struct {
	db         *gorm.DB
	repository repository.Repository
}

func NewContactDetailsService(DB *gorm.DB, repo repository.Repository) *ContactDeatailsService {
	return &ContactDeatailsService{
		db:         DB,
		repository: repo,
	}
}

func (service *ContactDeatailsService) CreateContactDetail(newContactDetail *contactdetail.ContactDetail) error {
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	if err := service.repository.Add(uow, newContactDetail); err != nil {
		return apperror.NewDatabaseError("Failed to create contact detail: " + err.Error())
	}

	uow.Commit()
	return nil
}
