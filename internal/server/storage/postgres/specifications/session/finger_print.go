package session

import (
	"fmt"
	"github.com/feiin/sqlstring"

	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

var _ ports.QuerySpecification = &FingerPrintSpecification{}

type FingerPrintSpecification struct {
	FingerPrint string
}

func (s *FingerPrintSpecification) Query() string {
	return fmt.Sprintf("finger_print = %v", sqlstring.Escape(s.FingerPrint))
}

func (s *FingerPrintSpecification) Value() []any {
	return []any{s.FingerPrint}
}
