package contact

import (
	"contact_app_mux_gorm_main/components/log"

	"github.com/jinzhu/gorm"
)

type ContactModuleConfig struct {
	DB *gorm.DB
}

func NewContactModuleConfig(db *gorm.DB) *ContactModuleConfig {
	return &ContactModuleConfig{
		DB: db,
	}
}

func (c *ContactModuleConfig) MigrateTables() {
	model := &Contact{}

	err := c.DB.AutoMigrate(model).Error
	if err != nil {
		log.NewLog().Print("Auto Migrating Contact ==> %s", err)
	}

	err = c.DB.Model(model).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.NewLog().Print("Foreign Key Constraints Of Contact ==> %s", err)
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

// type Contact struct {
// 	UserID         uint                           `json:"UserID"`
// 	FName          string                         `json:"FirstName" gorm:"not null;type:varchar(100)"`
// 	LName          string                         `json:"LastName" gorm:"not null;type:varchar(100)"`
// 	IsActive       bool                           `json:"IsActive" gorm:"type:boolean;default:true"`
// 	ContactDetails []*contactdetail.ContactDetail `gorm:"foreignKey:ContactID"`
// }
