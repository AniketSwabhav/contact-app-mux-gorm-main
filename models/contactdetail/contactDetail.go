package contactdetail

import "contact_app_mux_gorm_main/models"

type ContactDetail struct {
	models.Base
	ContactID      string      `json:"ContactID" gorm:"type:varchar(36);not null"`
	TypeOfContact  string      `json:"Type" gorm:"not null;type:varchar(100)"`
	ValueOfContact interface{} `json:"value" gorm:"not null;type:varchar(100)"`
}
