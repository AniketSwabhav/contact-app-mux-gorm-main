package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/contact/controller"
	"contact_app_mux_gorm_main/components/contact/service"
	"contact_app_mux_gorm_main/modules/repository"
)

func registerContactRoutes(appObj *app.App, repository repository.Repository) {

	defer appObj.WG.Done()

	contactService := service.NewContactService(appObj.DB, repository)

	contactController := controller.NewContactController(contactService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactController,
	})
}
