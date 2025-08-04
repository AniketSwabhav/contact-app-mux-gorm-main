package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/contact/controller"
	"contact_app_mux_gorm_main/components/contact/service"
)

func registerContactRoutes(appObj *app.App) {

	contactService := service.NewContactService(appObj.DB)

	contactController := controller.NewContactController(contactService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactController,
	})
}
