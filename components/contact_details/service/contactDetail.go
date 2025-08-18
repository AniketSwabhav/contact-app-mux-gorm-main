package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/models/contactdetail"
	"contact_app_mux_gorm_main/modules/repository"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
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

func (service *ContactDeatailsService) GetAllContactDetail(contactId uuid.UUID, allContactDetails *[]contactdetail.ContactDetailDTO, totalCount *int) error {

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	queryProcessors := []repository.QueryProcessor{
		repository.Filter("contact_id = ?  AND deleted_at is NULL", contactId),
	}

	err := service.repository.GetAll(uow, allContactDetails, queryProcessors...)
	if err != nil {
		return apperror.NewDatabaseError("Failed to get all contact-Details")
	}

	uow.Commit()
	return nil
}

func (service *ContactDeatailsService) GetContactDetailById(targetContactDetail *contactdetail.ContactDetail) error {
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()

	queryProcessors := []repository.QueryProcessor{
		repository.Filter("id = ? AND contact_id = ?", targetContactDetail.ID, targetContactDetail.ContactID),
	}

	err := service.repository.GetRecord(uow, targetContactDetail, queryProcessors...)
	if err != nil {
		return apperror.NewNotFoundError("Contact-Detail not found")
	}

	uow.Commit()
	return nil
}

func (service *ContactDeatailsService) UpdateContactDetailById(targetContactDetail *contactdetail.ContactDetail) error {
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	if err := service.repository.Update(uow, targetContactDetail); err != nil {
		return apperror.NewDatabaseError("Failed to update contact detail: " + err.Error())
	}

	uow.Commit()
	return nil
}

func (service *ContactDeatailsService) DeleteContactDetailById(contactDetailID, contactID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	var contactDetailToDelete contactdetail.ContactDetail

	queryProcessors := []repository.QueryProcessor{
		repository.Filter("id = ? AND contact_id = ?", contactDetailID, contactID),
	}

	err := service.repository.GetRecord(uow, &contactDetailToDelete, queryProcessors...)
	if err != nil {
		return apperror.NewNotFoundError("Contact not found")
	}

	err = service.repository.UpdateWithMap(uow, &contactDetailToDelete, map[string]interface{}{
		"DeletedAt": time.Now(),
	}, queryProcessors...)
	if err != nil {
		return apperror.NewDatabaseError("Failed to delete contact- detail:")
	}

	uow.Commit()
	return nil
}
