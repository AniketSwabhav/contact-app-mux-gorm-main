package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/credential/controller"
	"contact_app_mux_gorm_main/components/credential/service"
)

func registerCredentialRoutes(appObj *app.App) {

	credentialService := service.NewCredentialService(appObj.DB)

	credentialController := controller.NewCredentialController(credentialService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		credentialController,
	})
}
