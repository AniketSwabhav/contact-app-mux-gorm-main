package modules

import "contact_app_mux_gorm_main/app"

func RegisterModuleRoutes(app *app.App) {

	log := app.Log
	log.Print("============Registering-Module-Routes==============")

	registerUserRoutes(app)
	registerContactRoutes(app)
	// registerContactInfoRoutes(app)
	registerCredentialRoutes(app)
}
