package user

import (
	"database/sql"
	"fmt"

	"github.com/feiin/sqlstring"

	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &EmailSpecification{}

type EmailSpecification struct {
	sql.DB
	Email string
}

func (s *EmailSpecification) Query() string {
	return fmt.Sprintf("email = %v", sqlstring.Escape(s.Email))
}

func (s *EmailSpecification) Value() []any {
	return []any{s.Email}
}
