package credential

import (
	"contact_app_mux_gorm_main/models"
)

type Credentials struct {
	models.Base
	Email    string `json:"Email" gorm:"unique;not null;type:varchar(100)"`
	Password string `json:"Password" gorm:"not null;type:varchar(100)"`
}
