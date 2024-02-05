package user

import (
	"database/sql"
	"fmt"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ secondary_ports.QuerySpecification = &EmailSpecification{}

type EmailSpecification struct {
	sql.DB
	Email string
}

func (s *EmailSpecification) Query() string {
	return fmt.Sprintf("email = '%s'", s.Email)
}

func (s *EmailSpecification) Value() []any {
	return []any{s.Email}
}
