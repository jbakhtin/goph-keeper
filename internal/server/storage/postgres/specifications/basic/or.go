package basic

import (
	"fmt"
	"strings"

	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &WhereSpecification{}

type OrSpecification struct {
	Specifications []ports.QuerySpecification
}

func (s *OrSpecification) Query() string {
	var queries []string
	for _, specification := range s.Specifications {
		queries = append(queries, specification.Query())
	}

	query := strings.Join(queries, " OR ")

	return fmt.Sprintf("(%s)", query)
}

func (s *OrSpecification) Value() []any {
	var values []interface{}
	for _, specification := range s.Specifications {
		values = append(values, specification.Value()...)
	}
	return values
}
