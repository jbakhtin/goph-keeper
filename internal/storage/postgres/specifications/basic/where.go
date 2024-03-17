package basic

import (
	"fmt"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"
)

var _ queryspec.QuerySpecification = &WhereSpecification{}

type WhereSpecification struct {
	Specification queryspec.QuerySpecification
}

func (s *WhereSpecification) Query() string {
	return fmt.Sprintf("WHERE (%s)", s.Specification.Query())
}

func (s *WhereSpecification) Value() []any {
	return s.Specification.Value()
}
