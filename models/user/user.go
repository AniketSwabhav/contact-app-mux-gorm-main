package user

import (
	"contact_app_mux_gorm_main/models"
	"contact_app_mux_gorm_main/models/contact"
	"contact_app_mux_gorm_main/models/credential"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	models.Base
	FirstName     string                  `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName      string                  `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsAdmin       bool                    `json:"IsAdmin" gorm:"type:boolean;default:false"`
	IsActive      bool                    `json:"IsActive" gorm:"type:boolean;default:true"`
	CredentialsID uuid.UUID               `json:"CredentialID" gorm:"type:varchar(36)"`
	Credentials   *credential.Credentials `json:"-" gorm:"foreignKey:CredentialsID;references:ID"`
	Contacts      []contact.Contact       `json:"-" gorm:"foreignKey:UserID;references:UserID"`
}

// func CreateAdmin(fName, lName string, credential *credential.Credentials) *User {

// 	id := uuid.New()

// 	user := &User{
// 		UserID:      id.String(),
// 		FirstName:   fName,
// 		LastName:    lName,
// 		IsAdmin:     true,
// 		IsActive:    true,
// 		Credentials: credential,
// 	}

// 	return user
// }

// func CreateUser(fName, lName string, credential *credential.Credentials) *User {

// 	id := uuid.New()

// 	user := &User{
// 		UserID:      id.String(),
// 		FirstName:   fName,
// 		LastName:    lName,
// 		IsAdmin:     false,
// 		IsActive:    true,
// 		Credentials: credential,
// 	}

// 	return user
// }
