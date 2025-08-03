package user

import (
	"contact_app_mux_gorm_main/models/contact"
	"contact_app_mux_gorm_main/models/credential"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID        string                  `json:"UserID" gorm:"not null;unique;type:varchar(100);primaryKey"`
	FirstName     string                  `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName      string                  `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsAdmin       bool                    `json:"IsAdmin" gorm:"type:boolean;default:false"`
	IsActive      bool                    `json:"IsActive" gorm:"type:boolean;default:true"`
	CredentialsID uint                    `gorm:"not null;unique"`
	Credentials   *credential.Credentials `json:"Credentials" gorm:"foreignKey:CredentialsID;references:ID"`
	Contacts      []contact.Contact       `json:"Contacts" gorm:"foreignKey:UserID;references:UserID"`
}

func CreateAdmin(fName, lName string, credential *credential.Credentials) *User {

	id := uuid.New()

	user := &User{
		UserID:      id.String(),
		FirstName:   fName,
		LastName:    lName,
		IsAdmin:     true,
		IsActive:    true,
		Credentials: credential,
	}

	return user
}

func CreateUser(fName, lName string, credential *credential.Credentials) *User {

	id := uuid.New()

	user := &User{
		UserID:      id.String(),
		FirstName:   fName,
		LastName:    lName,
		IsAdmin:     false,
		IsActive:    true,
		Credentials: credential,
	}

	return user
}
