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
	log                  log.Log
	contactDetailService *service.ContactDeatailsService
}

func NewContacDetailsController(contactDetailsService *service.ContactDeatailsService, log log.Log) *ContactDetailsController {
	return &ContactDetailsController{
		log:                  log,
		contactDetailService: contactDetailsService,
	}
}

func (cd *ContactDetailsController) RegisterRoutes(router *mux.Router) {
	contactDetailsRouter := router.PathPrefix("/user/{userId}/contact/{contactId}").Subrouter()

	//Post
	contactDetailsRouter.HandleFunc("/contactDetail/", cd.createContactDetail).Methods(http.MethodPost)

	//Get
	contactDetailsRouter.HandleFunc("/contactDetail/", cd.getAllContactDetail).Methods(http.MethodGet)
	contactDetailsRouter.HandleFunc("/contactDetail/{contactDetailId}", cd.getContactDetailById).Methods(http.MethodGet)

	//updaet
	contactDetailsRouter.HandleFunc("/contactDetail/{contactDetailId}", cd.updateContactDetailById).Methods(http.MethodPut)

	//delete
	contactDetailsRouter.HandleFunc("/contactDetail/{contactDetailId}", cd.deleteContactDetailById).Methods(http.MethodDelete)

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

	fmt.Println(newContactDetail)

	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]
	contactIdFromURL := vars["contactId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid user ID format"))
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

	err = cd.contactDetailService.CreateContactDetail(&newContactDetail)
	if err != nil {
		cd.log.Print(err.Error())
		util.RespondError(w, apperror.NewHTTPError(err.Error()))
		return
	}

	util.RespondJSON(w, http.StatusCreated, newContactDetail)
}

func (cd *ContactDetailsController) getAllContactDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]
	contactIdFromURL := vars["contactId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid user ID format"))
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
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to get a contact for this user"))
		return
	}

	allContactDetails := &[]contactdetail.ContactDetailDTO{}
	var totalCount int

	err = cd.contactDetailService.GetAllContactDetail(contactUUID, allContactDetails, &totalCount)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, allContactDetails)
}

func (cd *ContactDetailsController) getContactDetailById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]
	contactIdFromURL := vars["contactId"]
	contactDetailIdFromURL := vars["contactDetailId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid user ID format"))
		return
	}

	contactUUID, err := uuid.FromString(contactIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact ID format"))
		return
	}

	contactDetailUUID, err := uuid.FromString(contactDetailIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact-Detail ID format"))
		return
	}

	claim := authorization.Claims{}

	err = authorization.ValidateToken(w, r, &claim)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	if claim.UserID != userUUID {
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to get all contact-detail for this user"))
		return
	}

	targetContactDetail := contactdetail.ContactDetail{}

	targetContactDetail.ContactID = contactUUID
	targetContactDetail.ID = contactDetailUUID

	err = cd.contactDetailService.GetContactDetailById(&targetContactDetail)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, targetContactDetail)
}

func (cd *ContactDetailsController) updateContactDetailById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdFromURL := vars["userId"]
	contactIdFromURL := vars["contactId"]
	contactDetailIdFromURL := vars["contactDetailId"]

	userUUID, err := uuid.FromString(userIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid user ID format"))
		return
	}

	contactUUID, err := uuid.FromString(contactIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact ID format"))
		return
	}

	contactDetailUUID, err := uuid.FromString(contactDetailIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact-Detail ID format"))
		return
	}

	claim := authorization.Claims{}

	err = authorization.ValidateToken(w, r, &claim)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("BAD_GATEWAY", err.Error()))
		return
	}

	if claim.UserID != userUUID {
		util.RespondError(w, apperror.NewAuthorizationError("You are not authorized to update contact-detail for this user"))
		return
	}

	updateContactDetail := contactdetail.ContactDetail{}
	err = util.UnmarshalJSON(r, &updateContactDetail)
	if err != nil {
		util.RespondError(w, apperror.NewHTTPError(err.Error()))
		return
	}

	updateContactDetail.ID = contactDetailUUID
	updateContactDetail.ContactID = contactUUID

	err = cd.contactDetailService.UpdateContactDetailById(&updateContactDetail)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, updateContactDetail)
}

func (cd *ContactDetailsController) deleteContactDetailById(w http.ResponseWriter, r *http.Request) {
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

	contactDetailIdFromURL := vars["contactDetailId"]

	contactDetailUUID, err := uuid.FromString(contactDetailIdFromURL)
	if err != nil {
		util.RespondError(w, apperror.NewValidationError("INVALID_UUID", "Invalid Contact-Detail ID format"))
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

	err = cd.contactDetailService.DeleteContactDetailById(contactDetailUUID, contactUUID)
	if err != nil {
		util.RespondError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Contact deleted successfully",
	})
}
