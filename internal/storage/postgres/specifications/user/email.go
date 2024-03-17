package user

import (
	"database/sql"
	"fmt"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/feiin/sqlstring"
)

var _ queryspec.QuerySpecification = &EmailSpecification{}

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
