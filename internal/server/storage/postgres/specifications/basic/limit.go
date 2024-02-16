package basic

import (
	"fmt"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/feiin/sqlstring"
)

var _ queryspec.QuerySpecification = &WhereSpecification{}

type LimitSpecification struct {
	Specification queryspec.QuerySpecification
	Limit         int
}

func (s *LimitSpecification) Query() string {
	return fmt.Sprintf("%s LIMIT %v", s.Specification.Query(), sqlstring.Escape(s.Limit))
}

func (s *LimitSpecification) Value() []any {
	return append(s.Specification.Value(), s.Limit)
}
