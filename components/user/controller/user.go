package controller

import (
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/components/user/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserController struct {
	log     log.Log
	service *service.UserService
}

func NewUserController(userService *service.UserService, log log.Log) *UserController {
	return &UserController{
		service: userService,
		log:     log,
	}
}

type UserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (c *UserController) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	// guardedRouter := userRouter.PathPrefix("/").Subrouter()
	unguardedRouter := userRouter.PathPrefix("/").Subrouter()

	unguardedRouter.HandleFunc("/register-admin", c.registerAdmin).Methods(http.MethodPost)
}

func (c *UserController) registerAdmin(w http.ResponseWriter, r *http.Request) {

	var userInput *UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := c.service.CreateAdmin(userInput.FirstName, userInput.LastName, userInput.Email, userInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}
