package session

import (
	"fmt"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/feiin/sqlstring"
)

var _ queryspec.QuerySpecification = &UserIDSpecification{}

type UserIDSpecification struct {
	UserID int
}

func (s *UserIDSpecification) Query() string {
	return fmt.Sprintf("user_id = %v", sqlstring.Escape(s.UserID))
}

func (s *UserIDSpecification) Value() []any {
	return []any{s.UserID}
}
