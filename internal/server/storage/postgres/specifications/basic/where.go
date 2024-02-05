package basic

import (
	"fmt"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ secondary_ports.QuerySpecification = &WhereSpecification{}

type WhereSpecification struct {
	Specification secondary_ports.QuerySpecification
}

func (s *WhereSpecification) Query() string {
	return fmt.Sprintf("WHERE (%s)", s.Specification.Query())
}

func (s *WhereSpecification) Value() []any {
	return s.Specification.Value()
}
