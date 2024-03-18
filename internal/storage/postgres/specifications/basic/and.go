package basic

import (
	"strings"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"
)

var _ queryspec.QuerySpecification = &AndSpecification{}

type AndSpecification struct {
	Specifications []queryspec.QuerySpecification
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
