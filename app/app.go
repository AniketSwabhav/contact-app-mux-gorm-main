package app

import (
	"contact_app_mux_gorm_main/components/log"
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	sync.Mutex
	Name   string
	Router *mux.Router
	DB     *gorm.DB
	Log    log.Log
	Server *http.Server
}

type Controller interface {
	RegisterRoutes(router *mux.Router)
}

type ModuleConfig interface {
	MigrateTables()
}

func NewApp(name string, db *gorm.DB, log log.Log) *App {
	return &App{
		Name: name,
		DB:   db,
		Log:  log,
	}
}

func NewDBConnection(log log.Log) *gorm.DB {
	const url = "root:12345@tcp(127.0.0.1:3306)/contact_app_db?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Print(err.Error())
		return nil
	}

	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(3 * time.Minute)

	db.LogMode(true)

	return db
}

func (a *App) Init() {
	a.initializeRouter()
	a.initializeServer()
}

func (a *App) initializeRouter() {
	a.Log.Print("Initializing " + a.Name + " Route")
	a.Router = mux.NewRouter().StrictSlash(true)
	a.Router = a.Router.PathPrefix("/api/v1/contact-app").Subrouter()
}

func (a *App) initializeServer() {
	headersOk := handlers.AllowCredentials()
	originsOk := handlers.AllowedOrigins([]string{
		os.Getenv("ORIGIN_ALLOWED"),
	})
	methodsOk := handlers.AllowedMethods([]string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
	})

	a.Server = &http.Server{
		Addr:         "localhost:4002",
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(a.Router),
	}
	a.Log.Print("Server Exposed On 4002")
}

func (a *App) StartServer() error {

	a.Log.Print("Server Time: ", time.Now())
	a.Log.Print("Server Running on port:4002")

	err := a.Server.ListenAndServe()
	if err != nil {
		a.Log.Print("Listen and serve error: ", err)
		return err
	}
	return nil
}

func (a *App) RegisterControllerRoutes(controllers []Controller) {

	a.Lock()
	defer a.Unlock()

	for _, controller := range controllers {
		controller.RegisterRoutes(a.Router.NewRoute().Subrouter())
	}

}

func (a *App) MigrateModuleTables(moduleConfigs []ModuleConfig) {

	a.Lock()
	defer a.Unlock()

	for _, moduleConfig := range moduleConfigs {
		moduleConfig.MigrateTables()
	}

}

func (app *App) Stop() {

	context, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	app.DB.Close()
	app.Log.Print("Db closed")

	err := app.Server.Shutdown(context)
	if err != nil {
		app.Log.Print("Failed to Stop Server")
		return
	}
	app.Log.Print("Server Shutdown")
}
