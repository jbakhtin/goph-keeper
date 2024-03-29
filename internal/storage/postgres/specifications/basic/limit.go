package basic

import (
	"fmt"

	"github.com/feiin/sqlstring"

	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &WhereSpecification{}

type LimitSpecification struct {
	Specification ports.QuerySpecification
	Limit         int
}

func (s *LimitSpecification) Query() string {
	return fmt.Sprintf("%s LIMIT %v", s.Specification.Query(), sqlstring.Escape(s.Limit))
}

func (s *LimitSpecification) Value() []any {
	return append(s.Specification.Value(), s.Limit)
}
