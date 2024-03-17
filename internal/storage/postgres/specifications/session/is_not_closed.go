package session

import ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"

var _ ports.QuerySpecification = &IsNotClosedSpecification{}

type IsNotClosedSpecification struct{}

func (s *IsNotClosedSpecification) Query() string {
	return "closed_at IS NULL"
}

func (s *IsNotClosedSpecification) Value() []any {
	return nil
}
