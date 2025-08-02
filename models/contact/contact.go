package contact

import (
	"contact_app_mux_gorm_main/models/contactdetail"
)

type Contact struct {
	ContactID      string                         `json:"ContactID" gorm:"primaryKey;type:varchar(100);not null;unique"`
	FirstName      string                         `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName       string                         `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsActive       bool                           `json:"IsActive" gorm:"type:boolean;default:true"`
	UserID         string                         `json:"UserID" gorm:"type:varchar(100);not null"`
	ContactDetails []*contactdetail.ContactDetail `gorm:"foreignKey:ContactID;references:ContactID"`
}
