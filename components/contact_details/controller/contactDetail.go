package controller

import (
	"contact_app_mux_gorm_main/components/apperror"
	"contact_app_mux_gorm_main/components/contact_details/service"
	"contact_app_mux_gorm_main/components/log"
	"contact_app_mux_gorm_main/components/security/middleware/authorization"
	"contact_app_mux_gorm_main/components/util"
	"contact_app_mux_gorm_main/models/contactdetail"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type ContactDetailsController struct {
	log            log.Log
	contactService *service.ContactDeatailsService
}

func NewContacDetailsController(contactDetailsService *service.ContactDeatailsService, log log.Log) *ContactDetailsController {
	return &ContactDetailsController{
		log:            log,
		contactService: contactDetailsService,
	}
}

func (cd *ContactDetailsController) RegisterRoutes(router *mux.Router) {
	contactDetailsRouter := router.PathPrefix("/user/{userId}/contact/{contactId}").Subrouter()

	contactDetailsRouter.HandleFunc("/", cd.createContactDetail).Methods(http.MethodPost)
}

func (cd *ContactDetailsController) createContactDetail(w http.ResponseWriter, r *http.Request) {
	newContactDetail := contactdetail.ContactDetail{}

	err := util.UnmarshalJSON(r, &newContactDetail)
	if err != nil {
		fmt.Println("==============================err from UnmarshalJSON==========================")
		cd.log.Print(err.Error())
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

	newContactDetail.ContactID = contactUUID

	err = cd.contactService.CreateContactDetail(&newContactDetail)
	if err != nil {
		cd.log.Print(err.Error())
		util.RespondError(w, apperror.NewHTTPError(err.Error()))
		return
	}

	util.RespondJSON(w, http.StatusCreated, newContactDetail)
}
