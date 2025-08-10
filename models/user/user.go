package user

import (
	"contact_app_mux_gorm_main/models"
	"contact_app_mux_gorm_main/models/contact"
	"contact_app_mux_gorm_main/models/credential"
)

type User struct {
	models.Base
	FirstName   string                  `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName    string                  `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsAdmin     bool                    `json:"IsAdmin" gorm:"type:boolean;default:false"`
	IsActive    bool                    `json:"IsActive" gorm:"type:boolean;default:true"`
	Credentials *credential.Credentials `json:"Credentials" gorm:"foreignKey:UserID;references:ID"`
}

type UserDTO struct {
	models.Base
	FirstName   string                  `json:"FirstName" gorm:"not null;type:varchar(100)"`
	LastName    string                  `json:"LastName" gorm:"not null;type:varchar(100)"`
	IsAdmin     bool                    `json:"IsAdmin" gorm:"type:boolean;default:false"`
	IsActive    bool                    `json:"IsActive" gorm:"type:boolean;default:true"`
	Credentials *credential.Credentials `json:"Credentials" gorm:"foreignKey:UserID;references:ID"`
	Contacts    []contact.Contact       `json:"Contacts" gorm:"foreignKey:UserID;references:ID"`
}

func (*UserDTO) TableName() string {
	return "users"
}
