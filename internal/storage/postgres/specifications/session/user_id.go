package session

import (
	"fmt"

	"github.com/feiin/sqlstring"

	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &UserIDSpecification{}

type UserIDSpecification struct {
	UserID int
}

func (s *UserIDSpecification) Query() string {
	return fmt.Sprintf("user_id = %v", sqlstring.Escape(s.UserID))
}

func (s *UserIDSpecification) Value() []any {
	return []any{s.UserID}
}
