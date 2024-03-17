package basic

import (
	"fmt"

	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &WhereSpecification{}

type WhereSpecification struct {
	Specification ports.QuerySpecification
}

func (s *WhereSpecification) Query() string {
	return fmt.Sprintf("WHERE (%s)", s.Specification.Query())
}

func (s *WhereSpecification) Value() []any {
	return s.Specification.Value()
}
