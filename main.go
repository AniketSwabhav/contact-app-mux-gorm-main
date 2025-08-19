package main

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/config"
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/docs"
	"contact_app_mux_gorm_main/modules"
	"contact_app_mux_gorm_main/modules/repository"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var environment = "local"

func main() {

	env := config.Environment(environment)

	log := log.GetLogger()
	log.Info("Starting main in ", env, ".")

	config.InitializeGlobalConfig(env)

	if env == config.Local {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", config.PORT.GetStringValue())
	}

	db := app.NewDBConnection(log)
	if db == nil {
		log.Fatalf("Db connection failed.")
	}
	defer func() {
		db.Close()
		log.Info("Db closed")
	}()

	var wg sync.WaitGroup
	var repository = repository.NewGormRepository()

	app := app.NewApp("Contact App", db, log, &wg, repository)
	app.Init()

	modules.RegisterModuleRoutes(app, repository)

	go func() {
		err := app.StartServer()
		if err != nil {
			app.Log.Print("Error in starting Server")
			stopApp(app)
		}
	}()

	app.Log.Print("Server Started")

	modules.Configure(app)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	stopApp(app)
}

func stopApp(app *app.App) {
	app.Stop()
	app.Log.Print("App stopped.")
	os.Exit(0)
}
