package session

import (
	"fmt"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ secondary_ports.QuerySpecification = &FingerPrintSpecification{}

type FingerPrintSpecification struct {
	FingerPrint string
}

func (s *FingerPrintSpecification) Query() string {
	return fmt.Sprintf("finger_print = %v", s.FingerPrint)
}

func (s *FingerPrintSpecification) Value() []any {
	return []any{s.FingerPrint}
}
