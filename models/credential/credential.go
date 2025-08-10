package credential

import (
	"contact_app_mux_gorm_main/models"

	uuid "github.com/satori/go.uuid"
)

type Credentials struct {
	models.Base
	Email    string    `json:"Email" gorm:"unique;not null;type:varchar(36)"`
	Password string    `json:"Password" gorm:"not null;type:varchar(100)"`
	UserID   uuid.UUID `json:"UserID" gorm:"not null;type:varchar(100)"`
}

type CredentialsDTO struct {
	models.Base
	Email    string    `json:"Email" gorm:"unique;not null;type:varchar(100)"`
	Password string    `json:"Password" gorm:"not null;type:varchar(100)"`
	UserID   uuid.UUID `json:"UserID" gorm:"not null;type:varchar(36)"`
}

func (*CredentialsDTO) TableName() string {
	return "credentials"
}
