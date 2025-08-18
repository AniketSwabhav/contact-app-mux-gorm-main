package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/models/contact"
	"contact_app_mux_gorm_main/modules/repository"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type ContactService struct {
	db         *gorm.DB
	repository repository.Repository
}

func NewContactService(DB *gorm.DB, repo repository.Repository) *ContactService {
	return &ContactService{
		db:         DB,
		repository: repo,
	}
}

func (service *ContactService) CreateContact(newContact *contact.Contact) error {

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	newContact.IsActive = true

	if err := service.repository.Add(uow, newContact); err != nil {
		return apperror.NewDatabaseError("Failed to create user: " + err.Error())
	}

	uow.Commit()
	return nil
}

func (service *ContactService) GetAllContacts(userID uuid.UUID, allContacts *[]contact.ContactDTO, totalCount *int) error {

	limit := 5
	offset := 0

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, repository.Filter("user_id = ?", userID))
	queryProcessors = append(queryProcessors, repository.PreloadAssociations([]string{"ContactDetails"}))
	queryProcessors = append(queryProcessors, repository.Paginate(limit, offset, totalCount))

	err := service.repository.GetAll(uow, allContacts, repository.CombineQueries(queryProcessors))
	if err != nil {
		return apperror.NewDatabaseError("Failed to get all contacts")
	}

	uow.Commit()
	return nil
}

func (service *ContactService) GetContactById(targetContact *contact.ContactDTO) error {
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()

	queryProcessors := []repository.QueryProcessor{
		repository.Filter("id = ? AND user_id = ?", targetContact.ID, targetContact.UserID),
		repository.PreloadAssociations([]string{"ContactDetails"}),
	}

	err := service.repository.GetRecord(uow, targetContact, queryProcessors...)
	if err != nil {
		return apperror.NewNotFoundError("Contact not found")
	}

	uow.Commit()
	return nil
}

func (service *ContactService) UpdateContactById(contactToUpdate *contact.Contact) error {

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	tempContact := contact.Contact{}

	queryProcessors := []repository.QueryProcessor{
		repository.Filter("id = ? AND user_id = ?", contactToUpdate.ID, contactToUpdate.UserID),
	}

	err := service.repository.GetRecord(uow, &tempContact, queryProcessors...)
	if err != nil {
		return err
	}

	err = service.repository.UpdateWithMap(uow, contactToUpdate, map[string]interface{}{
		"first_name": contactToUpdate.FirstName,
		"last_name":  contactToUpdate.LastName,
		"is_active":  contactToUpdate.IsActive,
	})

	uow.Commit()
	return nil
}

func (service *ContactService) DeleteContactById(contactID, userID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	var contactToDelete contact.Contact

	queryProcessors := []repository.QueryProcessor{
		repository.Filter("id = ? AND user_id = ?", contactID, userID),
	}

	err := service.repository.GetRecord(uow, &contactToDelete, queryProcessors...)
	if err != nil {
		return apperror.NewNotFoundError("Contact not found")
	}

	err = service.repository.UpdateWithMap(uow, &contactToDelete, map[string]interface{}{
		"DeletedAt": time.Now(),
	}, queryProcessors...)
	if err != nil {
		return apperror.NewDatabaseError("Failed to delete contact")
	}

	uow.Commit()
	return nil
}
