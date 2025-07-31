package contactdetail

type ContactDetail struct {
	ContactID      uint        `json:"ContactID"`
	TypeOfContact  string      `json:"Type" gorm:"not null;type:varchar(100)"`
	ValueOfContact interface{} `json:"value" gorm:"not null;type:varchar(100)"`
}
