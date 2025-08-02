package contactdetail

type ContactDetail struct {
	ContactDetailID string      `json:"ContactDetailID" gorm:"primaryKey;type:varchar(100);not null;unique"`
	ContactID       string      `json:"ContactID" gorm:"type:varchar(100);not null"`
	TypeOfContact   string      `json:"Type" gorm:"not null;type:varchar(100)"`
	ValueOfContact  interface{} `json:"value" gorm:"not null;type:varchar(100)"`
}
