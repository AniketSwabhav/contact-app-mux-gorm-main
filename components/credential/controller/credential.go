package controller

import (
	"contact_app_mux_gorm_main/components/credential/service"
	"contact_app_mux_gorm_main/components/log"

	"github.com/gorilla/mux"
)

type CredentialController struct {
	log     log.Log
	service *service.CredentialService
}

func NewCredentialController(credentialService *service.CredentialService, log log.Log) *CredentialController {
	return &CredentialController{
		log:     log,
		service: credentialService,
	}
}

func (c *CredentialController) RegisterRoutes(router *mux.Router) {}
