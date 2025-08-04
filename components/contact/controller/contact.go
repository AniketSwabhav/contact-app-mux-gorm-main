package controller

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/components/contact/service"
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"contact_app_mux_gorm_main/components/util"
	"contact_app_mux_gorm_main/models/contact"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ContactController struct {
	log     log.Log
	service *service.ContactService
}

func NewContactController(contactService *service.ContactService, log log.Log) *ContactController {
	return &ContactController{
		log:     log,
		service: contactService,
	}
}

func (c *ContactController) RegisterRoutes(router *mux.Router) {

	contactRouter := router.PathPrefix("/user/{userId}/contact").Subrouter()

	// POST
	contactRouter.HandleFunc("/", c.createContact).Methods(http.MethodPost)

	//Get
	contactRouter.HandleFunc("/", c.getAllContacts).Methods(http.MethodGet)
	contactRouter.HandleFunc("/{contactId}", c.getContactById).Methods(http.MethodGet)

	//Update
	contactRouter.HandleFunc("/{contactId}", c.updateContactById).Methods(http.MethodPut)

	contactRouter.Use(util.MiddlewareContact)

}

func (c *ContactController) createContact(w http.ResponseWriter, r *http.Request) {

	var contact *contact.Contact
	json.NewDecoder(r.Body).Decode(&contact)

	claim, err := authorization.ValidateToken(w, r)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	newContact, err := c.service.CreateContact(claim.UserID, contact)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, newContact)
}

func (c *ContactController) getAllContacts(w http.ResponseWriter, r *http.Request) {

	claim, err := authorization.ValidateToken(w, r)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	allContacts, err := c.service.GetAllContacts(claim.UserID)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, allContacts)
}

func (c *ContactController) getContactById(w http.ResponseWriter, r *http.Request) {

	claim, err := authorization.ValidateToken(w, r)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	vars := mux.Vars(r)
	contactIdFromURL := vars["contactId"]

	foundContact, err := c.service.GetContact(claim.UserID, contactIdFromURL)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, foundContact)
}

func (c *ContactController) updateContactById(w http.ResponseWriter, r *http.Request) {

	claim, err := authorization.ValidateToken(w, r)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	vars := mux.Vars(r)
	contactIdFromURL := vars["contactId"]

	var contact *contact.Contact
	json.NewDecoder(r.Body).Decode(&contact)

	updatedContact, err := c.service.UpdateContact(claim.UserID, contactIdFromURL, contact.FirstName, contact.LastName)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, updatedContact)
}
