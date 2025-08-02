package controller

import (
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"contact_app_mux_gorm_main/components/user/service"
	"contact_app_mux_gorm_main/models/credential"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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
	unguardedRouter.HandleFunc("/login", c.login).Methods(http.MethodPost)
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

func (c *UserController) login(w http.ResponseWriter, r *http.Request) {

	var userCredentials *credential.Credentials
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userCredentials.ValidateCredential()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	foundUser, err := c.service.Login(userCredentials.Email, userCredentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claim := &authorization.Claims{
		UserID:  foundUser.ID,
		IsAdmin: foundUser.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
		},
	}

	token, err := claim.Coder()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: token,
	})

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(foundUser)
}
