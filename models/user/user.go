package user

import (
	"contact_app_mux_gorm_main/models/contact"
	"contact_app_mux_gorm_main/models/credential"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName     string                  `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName      string                  `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsAdmin       bool                    `json:"IsAdmin" gorm:"type:boolean;default:false"`
	IsActive      bool                    `json:"IsActive" gorm:"type:boolean;default:true"`
	CredentialsID uint                    `json:"CredentialsID"`
	Credentials   *credential.Credentials `gorm:"foreignKey:CredentialsID"`
	Contacts      []contact.Contact       `gorm:"foreignKey:UserID"`
}

func CreateAdmin(fName, lName string, credential *credential.Credentials) *User {

	user := &User{
		FirstName:   fName,
		LastName:    lName,
		IsAdmin:     true,
		IsActive:    true,
		Credentials: credential,
	}

	return user
}
