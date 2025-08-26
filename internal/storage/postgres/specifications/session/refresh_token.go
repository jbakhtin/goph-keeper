package session

import (
	"fmt"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/feiin/sqlstring"
)

var _ queryspec.QuerySpecification = &RefreshTokenSpecification{}

type RefreshTokenSpecification struct {
	RefreshToken string
}

func (s *RefreshTokenSpecification) Query() string {
	return fmt.Sprintf("refresh_token = %v", sqlstring.Escape(s.RefreshToken))
}

func (s *RefreshTokenSpecification) Value() []any {
	return []any{s.RefreshToken}
}
