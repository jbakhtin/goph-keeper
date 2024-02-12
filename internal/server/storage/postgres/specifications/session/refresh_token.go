package session

import (
	"fmt"

	"github.com/feiin/sqlstring"

	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &RefreshTokenSpecification{}

type RefreshTokenSpecification struct {
	RefreshToken string
}

func (s *RefreshTokenSpecification) Query() string {
	return fmt.Sprintf("refresh_token = %v", sqlstring.Escape(s.RefreshToken))
}

func (s *RefreshTokenSpecification) Value() []any {
	return []any{s.RefreshToken}
}
