package session

import (
	"fmt"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/feiin/sqlstring"
)

var _ queryspec.QuerySpecification = &FingerPrintSpecification{}

type FingerPrintSpecification struct {
	FingerPrint string
}

func (s *FingerPrintSpecification) Query() string {
	return fmt.Sprintf("finger_print = %v", sqlstring.Escape(s.FingerPrint))
}

func (s *FingerPrintSpecification) Value() []any {
	return []any{s.FingerPrint}
}
