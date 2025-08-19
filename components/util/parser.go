package util

import (
	"contact_app_mux_gorm_main/components/apperror"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Parser struct {
	Params map[string]string
	Form   url.Values
}

func NewParser(r *http.Request) *Parser {
	r.ParseForm()
	return &Parser{
		Params: mux.Vars(r),
		Form:   r.Form,
	}
}

// GetUUID will get uuid from the given paramName in URL params.
func (p *Parser) GetUUID(paramName string) (uuid.UUID, error) {
	idString := p.Params[paramName]
	id, err := ParseUUID(idString)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func ParseUUID(input string) (uuid.UUID, error) {
	if len(input) == 0 {
		return uuid.Nil, apperror.NewValidationError("Empty ID", "ID cannot be empty")
	}
	id, err := uuid.FromString(input)
	if err != nil {
		return uuid.Nil, apperror.NewValidationError(input+": Invalid ID", "Invalid UUID format")
	}
	return id, nil
}
