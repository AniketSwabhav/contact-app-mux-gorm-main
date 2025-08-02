package credential

import (
	"contact_app_mux_gorm_main/components/log"

	"github.com/jinzhu/gorm"
)

type CredentialModuleConfig struct {
	DB *gorm.DB
}

func NewCredentialModuleConfig(db *gorm.DB) *CredentialModuleConfig {
	return &CredentialModuleConfig{
		DB: db,
	}
}

func (c *CredentialModuleConfig) MigrateTables() {

	model := &Credentials{}

	err := c.DB.AutoMigrate(model).Error
	if err != nil {
		log.NewLog().Print("Auto Migrating Credential ==> %s", err)
	}

}
