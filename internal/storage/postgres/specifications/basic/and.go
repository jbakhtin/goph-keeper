package basic

import (
	"strings"

	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &WhereSpecification{}

type AndSpecification struct {
	Specifications []ports.QuerySpecification
}

func (s *AndSpecification) Query() string {
	var queries []string
	for _, specification := range s.Specifications {
		queries = append(queries, specification.Query())
	}

	query := strings.Join(queries, " AND ")

	return query
}

func (s *AndSpecification) Value() []any {
	var values []interface{}
	for _, specification := range s.Specifications {
		values = append(values, specification.Value()...)
	}
	return values
}
