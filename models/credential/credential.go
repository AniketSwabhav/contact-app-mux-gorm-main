package credential

import (
	"github.com/jinzhu/gorm"
)

type Credentials struct {
	gorm.Model
	// CredentialID string `json:"CredentialID" gorm:"primary_key;type:varchar(100);not null;unique"`
	Email    string `json:"Email" gorm:"unique;not null;type:varchar(100)"`
	Password string `json:"Password" gorm:"not null;type:varchar(100)"`
}
