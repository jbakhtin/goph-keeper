package basic

import (
	"fmt"
	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &WhereSpecification{}

type LimitSpecification struct {
	Specification ports.QuerySpecification
	Limit         int
}

func (s *LimitSpecification) Query() string {
	return fmt.Sprintf("%s LIMIT %v", s.Specification.Query(), s.Limit)
}

func (s *LimitSpecification) Value() []any {
	return append(s.Specification.Value(), s.Limit)
}
