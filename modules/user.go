package modules

import (
	"contact_app_mux_gorm_main/app"
	credentialService "contact_app_mux_gorm_main/components/credential/service"
	"contact_app_mux_gorm_main/components/user/controller"
	userService "contact_app_mux_gorm_main/components/user/service"
	"contact_app_mux_gorm_main/modules/repository"
)

func registerUserRoutes(appObj *app.App, repository repository.Repository) {

	defer appObj.WG.Done()

	credentialService := credentialService.NewCredentialService(appObj.DB, repository)
	userService := userService.NewUserService(appObj.DB, repository)

	userController := controller.NewUserController(credentialService, userService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		userController,
	})
}
