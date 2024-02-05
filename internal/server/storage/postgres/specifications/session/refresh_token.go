package session

import (
	"fmt"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ secondary_ports.QuerySpecification = &RefreshTokenSpecification{}

type RefreshTokenSpecification struct {
	RefreshToken string
}

func (s *RefreshTokenSpecification) Query() string {
	return fmt.Sprintf("refresh_token = '%s'", s.RefreshToken)
}

func (s *RefreshTokenSpecification) Value() []any {
	return []any{s.RefreshToken}
}
