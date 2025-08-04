package service

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/models/contact"

	"github.com/jinzhu/gorm"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(DB *gorm.DB) *ContactService {
	return &ContactService{
		db: DB,
	}
}

func (c *ContactService) CreateContact(userID string, userContact *contact.Contact) (*contact.Contact, error) {

	newContact := contact.CreateContact(userContact.FirstName, userContact.LastName, userID)

	err := c.db.Create(newContact).Error
	if err != nil {
		return nil, apperror.NewDatabaseError("Failed to create contact")
	}

	return newContact, nil
}

func (c *ContactService) GetAllContacts(userID string) ([]contact.Contact, error) {

	allContacts := []contact.Contact{}
	err := c.db.Where("user_id = ?", userID).Find(&allContacts).Error
	if err != nil {
		return nil, apperror.NewDatabaseError("Failed to fetch contacts")
	}

	return allContacts, nil
}

func (c *ContactService) GetContact(userID string, contactID string) (*contact.Contact, error) {

	foundCountact := contact.Contact{}
	err := c.db.Where("user_id = ?", userID).Where("contact_id = ?", contactID).Find(&foundCountact).Error
	if err != nil {
		return nil, apperror.NewNotFoundError("Contact not found with given id")
	}

	return &foundCountact, nil
}

func (c *ContactService) UpdateContact(userID, contactID, firstName, lastName string) (*contact.Contact, error) {
	var foundContact contact.Contact

	err := c.db.Where("user_id = ?", userID).Where("contact_id = ?", contactID).First(&foundContact).Error
	if err != nil {
		return nil, apperror.NewNotFoundError("Contact not found with given ID")
	}

	foundContact.FirstName = firstName
	foundContact.LastName = lastName

	err = c.db.Save(&foundContact).Error
	if err != nil {
		return nil, apperror.NewDatabaseError("Failed to update contact")
	}

	return &foundContact, nil
}
