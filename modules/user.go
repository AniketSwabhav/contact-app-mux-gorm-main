package modules

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/user/controller"
	"contact_app_mux_gorm_main/components/user/service"
)

func registerUserRoutes(appObj *app.App) {

	userService := service.NewUserService(appObj.DB)

	userController := controller.NewUserController(userService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		userController,
	})
}
