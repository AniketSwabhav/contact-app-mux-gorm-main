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

	err = c.DB.Model(model).AddForeignKey("id", "users(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		log.NewLog().Print("Foreign Key Constraints Of Contact ==> %s", err)
	}
}
