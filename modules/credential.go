package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/credential/controller"
	"contact_app_mux_gorm_main/components/credential/service"
	"contact_app_mux_gorm_main/modules/repository"
)

func registerCredentialRoutes(appObj *app.App, repository repository.Repository) {

	credentialService := service.NewCredentialService(appObj.DB, repository)

	credentialController := controller.NewCredentialController(credentialService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		credentialController,
	})
}
