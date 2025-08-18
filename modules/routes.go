package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/modules/repository"
)

func RegisterModuleRoutes(app *app.App, repository repository.Repository) {

	log := app.Log
	log.Print("============Registering-Module-Routes==============")

	app.WG.Add(5)
	registerUserRoutes(app, repository)
	registerContactRoutes(app, repository)
	registerContactDetailRoutes(app, repository)
	registerCredentialRoutes(app, repository)

	app.WG.Done()
}
