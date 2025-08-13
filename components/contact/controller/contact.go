package controller

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/components/contact/service"
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"contact_app_mux_gorm_main/components/util"
	"contact_app_mux_gorm_main/models/contact"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type ContactController struct {
	log            log.Log
	contactService *service.ContactService
}

func NewContactController(contactService *service.ContactService, log log.Log) *ContactController {
	return &ContactController{
		log:            log,
		contactService: contactService,
	}
}

func (c *ContactController) RegisterRoutes(router *mux.Router) {

	contactRouter := router.PathPrefix("/user/{userId}/contact").Subrouter()

	// // POST
	contactRouter.HandleFunc("/", c.createContact).Methods(http.MethodPost)

	// //Get
	contactRouter.HandleFunc("/", c.getAllContacts).Methods(http.MethodGet)
	contactRouter.HandleFunc("/{contactId}", c.getContactById).Methods(http.MethodGet)

	// //Update
	contactRouter.HandleFunc("/{contactId}", c.updateContactById).Methods(http.MethodPut)

	//Delete
	contactRouter.HandleFunc("/{contactId}", c.deleteContactById).Methods(http.MethodDelete)

	contactRouter.Use(util.MiddlewareContact)

}

func (c *ContactController) createContact(w http.ResponseWriter, r *http.Request) {

	newContact := contact.Contact{}

	err := util.UnmarshalJSON(r, &newContact)
	if err != nil {
		fmt.Println("==============================err from UnmarshalJSON==========================")
		c.log.Print(err.Error())
		util.RespondError(w, apperror.NewHTTPError(err.Error()))
		return
	}

	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact ID format"))
		return
	}

	claim := authorization.Claims{}

	err = authorization.ValidateToken(w, r, &claim)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	if claim.UserID != userUUID {
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to create a contact for this user"))
		return
	}

	newContact.UserID = userUUID

	err = c.contactService.CreateContact(&newContact)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, newContact)
}

func (c *ContactController) getAllContacts(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact ID format"))
		return
	}

	claim := authorization.Claims{}

	err = authorization.ValidateToken(w, r, &claim)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	if claim.UserID != userUUID {
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to create a contact for this user"))
		return
	}

	allContacts := &[]contact.ContactDTO{}
	var totalCount int

	err = c.contactService.GetAllContacts(claim.UserID, allContacts, &totalCount)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	// util.RespondJSON(w, http.StatusOK, allContacts)
	util.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allContacts)
}

func (c *ContactController) getContactById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]
	contactIdFromURL := vars["contactId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_USER_ID", "Invalid user ID format"))
		return
	}

	contactUUID, err := uuid.FromString(contactIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_CONTACT_ID", "Invalid contact ID format"))
		return
	}

	claim := authorization.Claims{}
	err = authorization.ValidateToken(w, r, &claim)
	if err != nil {
		util.RespondError(w, apperror.NewAuthorizationError("Invalid token"))
		return
	}

	if claim.UserID != userUUID {
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to view this contact"))
		return
	}

	targetContact := contact.ContactDTO{}

	targetContact.ID = contactUUID
	targetContact.UserID = userUUID

	err = c.contactService.GetContactById(&targetContact)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, targetContact)
}

func (c *ContactController) updateContactById(w http.ResponseWriter, r *http.Request) {
	contactToBeUpdated := contact.Contact{}

	err := util.UnmarshalJSON(r, &contactToBeUpdated)
	if err != nil {
		fmt.Println("==============================err from UnmarshalJSON==========================")
		c.log.Print(err.Error())
		util.RespondError(w, apperror.NewHTTPError(err.Error()))
		return
	}

	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]
	contactIdFromURL := vars["contactId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact ID format"))
		return
	}

	contactUUID, err := uuid.FromString(contactIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_CONTACT_ID", "Invalid contact ID format"))
		return
	}

	claim := authorization.Claims{}

	err = authorization.ValidateToken(w, r, &claim)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	if claim.UserID != userUUID {
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to create a contact for this user"))
		return
	}
	contactToBeUpdated.UserID = userUUID
	contactToBeUpdated.ID = contactUUID

	err = c.contactService.UpdateContactById(&contactToBeUpdated)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, contactToBeUpdated)

}

func (c *ContactController) deleteContactById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]
	contactIdFromURL := vars["contactId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_USER_ID", "Invalid user ID format"))
		return
	}

	contactUUID, err := uuid.FromString(contactIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_CONTACT_ID", "Invalid contact ID format"))
		return
	}

	claim := authorization.Claims{}
	err = authorization.ValidateToken(w, r, &claim)
	if err != nil {
		util.RespondError(w, apperror.NewAuthorizationError("Invalid token"))
		return
	}

	if claim.UserID != userUUID {
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to delete this contact"))
		return
	}

	err = c.contactService.DeleteContactById(contactUUID, userUUID)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Contact deleted successfully",
	})
}
