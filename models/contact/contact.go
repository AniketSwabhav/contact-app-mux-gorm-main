package contact

import "contact_app_mux_gorm_main/models/contactdetail"

type Contact struct {
	UserID         uint                           `json:"UserID"`
	FirstName      string                         `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName       string                         `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsActive       bool                           `json:"IsActive" gorm:"type:boolean;default:true"`
	ContactDetails []*contactdetail.ContactDetail `gorm:"foreignKey:ContactID"`
}
