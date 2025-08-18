package contactdetail

import (
	"contact_app_mux_gorm_main/models"

	uuid "github.com/satori/go.uuid"
)

type ContactDetail struct {
	models.Base
	TypeOfContact  string      `json:"Type" gorm:"not null;type:varchar(100)"`
	ValueOfContact interface{} `json:"value" gorm:"not null;type:varchar(100)"`
	ContactID      uuid.UUID   `json:"ContactID" gorm:"type:varchar(36);not null"`
}

type ContactDetailDTO struct {
	// models.Base
	TypeOfContact  string      `json:"Type" gorm:"not null;type:varchar(100)"`
	ValueOfContact interface{} `json:"value" gorm:"not null;type:varchar(100)"`
	ContactID      uuid.UUID   `json:"ContactID" gorm:"type:varchar(36);not null"`
}

func (*ContactDetailDTO) TableName() string {
	return "contact_details"
}
