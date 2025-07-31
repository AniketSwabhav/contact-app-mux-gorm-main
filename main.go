package main

import (
	"contact_app_mux_gorm_main/app"
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/modules"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log := log.NewLog()
	db := app.NewDBConnection(*log)

	if db == nil {
		log.Print("DB connection falied")
	}
	defer func() {
		db.Close()
		log.Print("DB connection closed")
	}()

	app := app.NewApp("Contact App", db, *log)

	app.Init()

	modules.RegisterModuleRoutes(app)

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
