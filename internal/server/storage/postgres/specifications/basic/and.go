package basic

import (
	"fmt"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"strings"
)

var _ secondary_ports.QuerySpecification = &WhereSpecification{}

type AndSpecification struct {
	Specifications []secondary_ports.QuerySpecification
}

func (s *AndSpecification) Query() string {
	var queries []string
	for _, specification := range s.Specifications {
		queries = append(queries, specification.Query())
	}

	query := strings.Join(queries, " AND ")

	return fmt.Sprintf("%s", query)
}

func (s *AndSpecification) Value() []any {
	var values []interface{}
	for _, specification := range s.Specifications {
		values = append(values, specification.Value()...)
	}
	return values
}
