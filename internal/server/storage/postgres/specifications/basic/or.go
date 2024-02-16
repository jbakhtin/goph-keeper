package basic

import (
	"fmt"
	"strings"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"
)

var _ queryspec.QuerySpecification = &OrSpecification{}

type OrSpecification struct {
	Specifications []queryspec.QuerySpecification
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
