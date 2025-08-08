package controller

import (
	"contact_app_mux_gorm_main/components/apperror"
	credentialService "contact_app_mux_gorm_main/components/credential/service"
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	userService "contact_app_mux_gorm_main/components/user/service"
	"contact_app_mux_gorm_main/components/util"
	"contact_app_mux_gorm_main/models/credential"
	"contact_app_mux_gorm_main/models/user"
	"net/http"

	"github.com/gorilla/mux"
)

type UserController struct {
	log               log.Log
	UserService       *userService.UserService
	CredentialService *credentialService.CredentialService
}

func NewUserController(credentialService *credentialService.CredentialService, userService *userService.UserService, log log.Log) *UserController {
	return &UserController{
		CredentialService: credentialService,
		UserService:       userService,
		log:               log,
	}
}

// type UserInput struct {
// 	FirstName string
// 	LastName  string
// 	Email     string
// 	Password  string
// }

func (c *UserController) RegisterRoutes(router *mux.Router) {

	userRouter := router.PathPrefix("/user").Subrouter()
	guardedRouter := userRouter.PathPrefix("/").Subrouter()
	unguardedRouter := userRouter.PathPrefix("/").Subrouter()

	//Post
	unguardedRouter.HandleFunc("/login", c.login).Methods(http.MethodPost)

	unguardedRouter.HandleFunc("/register-admin", c.registerAdmin).Methods(http.MethodPost)
	guardedRouter.HandleFunc("/register", c.registerUser).Methods(http.MethodPost)

	// //Get
	// guardedRouter.HandleFunc("/", c.getAllUsers).Methods(http.MethodGet)
	// guardedRouter.HandleFunc("/get5", c.getUsersPaginated).Methods(http.MethodGet)
	// guardedRouter.HandleFunc("/{id}", c.getUserById).Methods(http.MethodGet)

	// //Update
	// guardedRouter.HandleFunc("/{id}", c.updateUserById).Methods(http.MethodPut)

	// //Delete
	// guardedRouter.HandleFunc("/{id}", c.deleteUserById).Methods(http.MethodDelete)

	guardedRouter.Use(util.MiddlewareAdmin)
}

func (c *UserController) registerAdmin(w http.ResponseWriter, r *http.Request) {

	newUser := user.User{}
	err := util.UnmarshalJSON(r, &newUser)
	if err != nil {
		util.RespondError(w, apperror.NewInvalidJSONError("Invalid JSON input"))
		return
	}

	err = c.UserService.CreateAdmin(&newUser)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, newUser)
}

func (c *UserController) registerUser(w http.ResponseWriter, r *http.Request) {

	newUser := user.User{}
	err := util.UnmarshalJSON(r, &newUser)
	if err != nil {
		util.RespondError(w, apperror.NewInvalidJSONError("Invalid JSON input"))
		return
	}

	err = c.UserService.CreateUser(&newUser)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, newUser)
}

func (c *UserController) login(w http.ResponseWriter, r *http.Request) {

	userCredentials := credential.Credentials{}
	err := util.UnmarshalJSON(r, &userCredentials)
	if err != nil {
		util.RespondError(w, apperror.NewInvalidJSONError("Invalid JSON input"))
		return
	}

	claim := authorization.Claims{}
	err = c.UserService.Login(&userCredentials, &claim)
	if err != nil {
		util.RespondError(w, err)
		return
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

	util.RespondJSON(w, http.StatusAccepted, map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}

// func (c *UserController) getAllUsers(w http.ResponseWriter, _ *http.Request) {

// 	allUsers, err := c.service.GetAllUsers()
// 	if err != nil {
// 		util.RespondError(w, err)
// 		return
// 	}

// 	util.RespondJSON(w, http.StatusAccepted, allUsers)
// }

// func (c *UserController) getUsersPaginated(w http.ResponseWriter, r *http.Request) {

// 	query := r.URL.Query()

// 	pageStr := query.Get("page")
// 	pageSizeStr := query.Get("limit")

// 	page := 1
// 	pageSize := 5

// 	if pageStr != "" {
// 		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
// 			page = parsedPage
// 		}
// 	}

// 	if pageSizeStr != "" {
// 		if parsedLimit, err := strconv.Atoi(pageSizeStr); err == nil && parsedLimit > 0 {
// 			pageSize = parsedLimit
// 		}
// 	}

// 	users, err := c.service.GetUsersPaginated(page, pageSize)
// 	if err != nil {
// 		util.RespondError(w, err)
// 		return
// 	}

// 	util.RespondJSON(w, http.StatusAccepted, users)
// }

// func (c *UserController) getUserById(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	userIdFromURL := vars["id"]

// 	foundUser, err := c.service.Get(userIdFromURL)
// 	if err != nil {
// 		util.RespondError(w, err)
// 		return
// 	}

// 	util.RespondJSON(w, http.StatusOK, foundUser)
// }

// func (c *UserController) updateUserById(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	userIdFromURL := vars["id"]

// 	var user *UserInput
// 	json.NewDecoder(r.Body).Decode(&user)
// 	targetUser, err := c.service.Update(userIdFromURL, user.FirstName, user.LastName, user.Email)
// 	if err != nil {
// 		util.RespondError(w, err)
// 		return
// 	}

// 	util.RespondJSON(w, http.StatusOK, targetUser)
// }

// func (c *UserController) deleteUserById(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userIdFromURL := vars["id"]

// 	err := c.service.Delete(userIdFromURL)
// 	if err != nil {
// 		util.RespondError(w, err)
// 		return
// 	}

// 	util.RespondJSON(w, http.StatusOK, map[string]string{
// 		"message": "User deleted successfully",
// 	})
// }
