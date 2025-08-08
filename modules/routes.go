package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/modules/repository"
)

func RegisterModuleRoutes(app *app.App, repository repository.Repository) {

	log := app.Log
	log.Print("============Registering-Module-Routes==============")

	registerUserRoutes(app, repository)
	// registerContactRoutes(app)
	// registerContactInfoRoutes(app)
	registerCredentialRoutes(app, repository)
}
