package user

import (
	"contact_app_mux_gorm_main/components/log"

	"github.com/jinzhu/gorm"
)

type UserModuleConfig struct {
	DB *gorm.DB
}

func NewUserModuleConfig(db *gorm.DB) *UserModuleConfig {
	return &UserModuleConfig{
		DB: db,
	}
}

func (u *UserModuleConfig) MigrateTables() {

	model := &User{}

	err := u.DB.AutoMigrate(model).Error
	if err != nil {
		log.NewLog().Print("Auto Migrating User ==> %s", err)
	}
}

// type User struct {
// 	gorm.Model
// 	FirstName     string                  `json:"FirstName" gorm:"not null;type:varchar(100)"`
// 	LastName      string                  `json:"LastName" gorm:"not null;type:varchar(100)"`
// 	IsAdmin       bool                    `json:"IsAdmin" gorm:"type:boolean;default:false"`
// 	IsActive      bool                    `json:"IsActive" gorm:"type:boolean;default:true"`
// 	CredentialsID uint                    `json:"CredentialsID"`
// 	Credentials   *credential.Credentials `gorm:"foreignKey:CredentialsID"`
// 	Contacts      []contact.Contact       `gorm:"foreignKey:UserID"`
// }
