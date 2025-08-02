package contactdetail

import (
	"contact_app_mux_gorm_main/components/log"

	"github.com/jinzhu/gorm"
)

type ContactInfoModuleConfig struct {
	DB *gorm.DB
}

func NewContactInfoModuleConfig(db *gorm.DB) *ContactInfoModuleConfig {
	return &ContactInfoModuleConfig{
		DB: db,
	}
}

func (c *ContactInfoModuleConfig) MigrateTables() {

	model := &ContactDetail{}

	err := c.DB.AutoMigrate(model).Error
	if err != nil {
		log.NewLog().Print("Auto Migrating ContactInfo ==> %s", err)
	}

	err = c.DB.Model(model).AddForeignKey("contact_id", "contacts(contact_id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.NewLog().Print("Foreign Key Constraints Of ContactDetail ==> %s", err)
	}
}
