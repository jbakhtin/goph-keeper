package session

import (
	"github.com/jbakhtin/goph-keeper/pkg/queryspec"
)

var _ queryspec.QuerySpecification = &IsNotClosedSpecification{}

type IsNotClosedSpecification struct{}

func (s *IsNotClosedSpecification) Query() string {
	return "closed_at IS NULL"
}

func (s *IsNotClosedSpecification) Value() []any {
	return nil
}
