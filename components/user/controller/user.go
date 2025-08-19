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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type UserController struct {
	log               log.Logger
	UserService       *userService.UserService
	CredentialService *credentialService.CredentialService
}

func NewUserController(credentialService *credentialService.CredentialService, userService *userService.UserService, log log.Logger) *UserController {
	return &UserController{
		CredentialService: credentialService,
		UserService:       userService,
		log:               log,
	}
}

func (c *UserController) RegisterRoutes(router *mux.Router) {

	userRouter := router.PathPrefix("/user").Subrouter()
	guardedRouter := userRouter.PathPrefix("/").Subrouter()
	unguardedRouter := userRouter.PathPrefix("/").Subrouter()

	//Post
	unguardedRouter.HandleFunc("/login", c.login).Methods(http.MethodPost)

	unguardedRouter.HandleFunc("/register-admin", c.registerAdmin).Methods(http.MethodPost)
	guardedRouter.HandleFunc("/register", c.registerUser).Methods(http.MethodPost)

	// //Get
	guardedRouter.HandleFunc("/", c.getAllUsers).Methods(http.MethodGet)
	guardedRouter.HandleFunc("/{id}", c.getUserById).Methods(http.MethodGet)

	// //Update
	guardedRouter.HandleFunc("/{id}", c.updateUserById).Methods(http.MethodPut)

	// //Delete
	guardedRouter.HandleFunc("/{id}", c.deleteUserById).Methods(http.MethodDelete)

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

	// http.SetCookie(w, &http.Cookie{
	// 	Name:  "auth",
	// 	Value: token,
	// })

	util.RespondJSON(w, http.StatusAccepted, map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}

func (controller *UserController) getAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers := &[]user.UserDTO{}
	var totalCount int
	err := controller.UserService.GetAllUsers(allUsers, &totalCount)
	if err != nil {
		controller.log.Print(err.Error())
		util.RespondError(w, err)
		return
	}
	util.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allUsers)
}

func (c *UserController) getUserById(w http.ResponseWriter, r *http.Request) {

	var targetUser = &user.UserDTO{}

	vars := mux.Vars(r)
	userIdFromURL := vars["id"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid user ID format"))
		return
	}

	targetUser.ID = userUUID

	err = c.UserService.GetUserByID(targetUser)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, targetUser)
}

func (c *UserController) updateUserById(w http.ResponseWriter, r *http.Request) {

	var userToUpdate = user.User{}

	err := util.UnmarshalJSON(r, &userToUpdate)
	if err != nil {
		fmt.Println("==============================err from UnmarshalJSON==========================")
		c.log.Print(err.Error())
		util.RespondError(w, apperror.NewHTTPError(err.Error()))
		return
	}
	vars := mux.Vars(r)
	userIdFromURL := vars["id"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid user ID format"))
		return
	}

	userToUpdate.ID = userUUID

	err = c.UserService.UpdateUser(&userToUpdate)
	if err != nil {
		c.log.Print(err.Error())
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, userToUpdate)
}

func (c *UserController) deleteUserById(w http.ResponseWriter, r *http.Request) {

	userToDelete := user.User{}

	vars := mux.Vars(r)
	userIdFromURL := vars["id"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid user ID format"))
		return
	}

	userToDelete.ID = userUUID

	err = c.UserService.Delete(&userToDelete)
	if err != nil {
		c.log.Print(err.Error())
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}
