package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/models/contact"
	"contact_app_mux_gorm_main/models/contactdetail"
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/models/user"
)

func Configure(appObj *app.App) {

	appObj.Log.Print("============Configuring-Module-Configs==============")

	userModule := user.NewUserModuleConfig(appObj.DB)
	contactModule := contact.NewContactModuleConfig(appObj.DB)
	contactInfoModule := contactdetail.NewContactInfoModuleConfig(appObj.DB)
	credentialModule := credential.NewCredentialModuleConfig(appObj.DB)

	appObj.MigrateModuleTables([]app.ModuleConfig{userModule, contactModule, contactInfoModule, credentialModule})

}
