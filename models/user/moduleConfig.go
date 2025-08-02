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

	// err = u.DB.Model(model).AddForeignKey("credentials_id", "credentials(credential_id)", "CASCADE", "CASCADE").Error
	// if err != nil {
	// 	log.NewLog().Print("Foreign Key Constraint User -> Credential ==> %s", err)
	// }
}
