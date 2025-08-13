package contact

import (
	"contact_app_mux_gorm_main/models"
	"contact_app_mux_gorm_main/models/contactdetail"

	uuid "github.com/satori/go.uuid"
)

type Contact struct {
	models.Base
	FirstName string    `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName  string    `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsActive  bool      `json:"IsActive" gorm:"type:boolean;default:true"`
	UserID    uuid.UUID `json:"UserID" gorm:"type:varchar(36);not null"`
}

type ContactDTO struct {
	models.Base
	FirstName      string                         `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName       string                         `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsActive       bool                           `json:"IsActive" gorm:"type:boolean;default:true"`
	UserID         uuid.UUID                      `json:"UserID" gorm:"type:varchar(36);not null"`
	ContactDetails []*contactdetail.ContactDetail `gorm:"foreignKey:ContactID;references:ID"`
}

func (*ContactDTO) TableName() string {
	return "contacts"
}
