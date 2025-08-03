package controller

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"contact_app_mux_gorm_main/components/user/service"
	"contact_app_mux_gorm_main/components/util"
	"contact_app_mux_gorm_main/models/credential"
	"encoding/json"
	"net/http"
	"strconv"
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
	guardedRouter := userRouter.PathPrefix("/").Subrouter()
	unguardedRouter := userRouter.PathPrefix("/").Subrouter()

	//Post
	unguardedRouter.HandleFunc("/login", c.login).Methods(http.MethodPost)

	unguardedRouter.HandleFunc("/register-admin", c.registerAdmin).Methods(http.MethodPost)
	guardedRouter.HandleFunc("/register", c.registerUser).Methods(http.MethodPost)

	//Get
	guardedRouter.HandleFunc("/", c.getAllUsers).Methods(http.MethodGet)
	guardedRouter.HandleFunc("/get5", c.getUsersPaginated).Methods(http.MethodGet)
	guardedRouter.HandleFunc("/{id}", c.getUserById).Methods(http.MethodGet)

	//Update
	guardedRouter.HandleFunc("/{id}", c.updateUserById).Methods(http.MethodPut)

	//Delete
	guardedRouter.HandleFunc("/{id}", c.deleteUserById).Methods(http.MethodDelete)

	guardedRouter.Use(util.MiddlewareAdmin)
}

func (c *UserController) registerAdmin(w http.ResponseWriter, r *http.Request) {

	var userInput *UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		util.RespondError(w, apperror.NewInvalidJSONError("Invalid JSON input"))
		return
	}

	createdUser, err := c.service.CreateAdmin(userInput.FirstName, userInput.LastName, userInput.Email, userInput.Password)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, createdUser)
}

func (c *UserController) registerUser(w http.ResponseWriter, r *http.Request) {

	var userInput *UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		util.RespondError(w, apperror.NewInvalidJSONError("Invalid JSON input"))
		return
	}

	createdUser, err := c.service.CreateUser(userInput.FirstName, userInput.LastName, userInput.Email, userInput.Password)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, createdUser)
}

func (c *UserController) login(w http.ResponseWriter, r *http.Request) {

	var userCredentials *credential.Credentials
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		util.RespondError(w, apperror.NewInvalidJSONError("Invalid JSON input"))
		return
	}

	err = userCredentials.ValidateCredential()
	if err != nil {
		util.RespondError(w, apperror.NewInValidTokenError(err.Error()))
		return
	}

	foundUser, err := c.service.Login(userCredentials.Email, userCredentials.Password)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	claim := &authorization.Claims{
		UserID:   foundUser.ID,
		IsAdmin:  foundUser.IsAdmin,
		IsActive: foundUser.IsActive,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
		},
	}

	token, err := claim.Coder()
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("TOKEN_CREATION_FAILED", "Could not create token"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: token,
	})

	util.RespondJSON(w, http.StatusAccepted, foundUser)
}

func (c *UserController) getAllUsers(w http.ResponseWriter, _ *http.Request) {

	allUsers, err := c.service.GetAllUsers()
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusAccepted, allUsers)
}

func (c *UserController) getUsersPaginated(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	pageStr := query.Get("page")
	pageSizeStr := query.Get("limit")

	page := 1
	pageSize := 5

	if pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if pageSizeStr != "" {
		if parsedLimit, err := strconv.Atoi(pageSizeStr); err == nil && parsedLimit > 0 {
			pageSize = parsedLimit
		}
	}

	users, err := c.service.GetUsersPaginated(page, pageSize)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusAccepted, users)
}

func (c *UserController) getUserById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIdFromURL := vars["id"]

	foundUser, err := c.service.Get(userIdFromURL)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, foundUser)
}

func (c *UserController) updateUserById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIdFromURL := vars["id"]

	var user *UserInput
	json.NewDecoder(r.Body).Decode(&user)
	targetUser, err := c.service.Update(userIdFromURL, user.FirstName, user.LastName, user.Email)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, targetUser)
}

func (c *UserController) deleteUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdFromURL := vars["id"]

	err := c.service.Delete(userIdFromURL)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}
