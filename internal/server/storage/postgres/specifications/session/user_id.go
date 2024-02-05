package session

import (
	"fmt"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ secondary_ports.QuerySpecification = &UserIDSpecification{}

type UserIDSpecification struct {
	UserID int
}

func (s *UserIDSpecification) Query() string {
	return fmt.Sprintf("user_id = %v", s.UserID)
}

func (s *UserIDSpecification) Value() []any {
	return []any{s.UserID}
}
