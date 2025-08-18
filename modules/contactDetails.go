package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/contact_details/controller"
	"contact_app_mux_gorm_main/components/contact_details/service"
	"contact_app_mux_gorm_main/modules/repository"
)

func registerContactDetailRoutes(appObj *app.App, repository repository.Repository) {

	defer appObj.WG.Done()

	contactDetailsService := service.NewContactDetailsService(appObj.DB, repository)

	contactDetailsController := controller.NewContacDetailsController(contactDetailsService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactDetailsController,
	})
}
