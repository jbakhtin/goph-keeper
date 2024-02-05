package session

import secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"

var _ secondary_ports.QuerySpecification = &IsNotClosedSpecification{}

type IsNotClosedSpecification struct{}

func (s *IsNotClosedSpecification) Query() string {
	return "closed_at IS NULL"
}

func (s *IsNotClosedSpecification) Value() []any {
	return nil
}
