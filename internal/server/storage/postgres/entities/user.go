package entities

import (
	"encoding/json"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"reflect"
	"time"
)

type FingerPrint map[any]any

func (fp *FingerPrint) Scan(value any) error {
	switch vt := value.(type) {
	case []byte:
		return json.Unmarshal(vt, &fp)
	default:
		return errors.New(fmt.Sprintf("can not convert %s to %s type", reflect.TypeOf(vt), "FingerPrint"))
	}
}

type User struct {
	ID int
	UUID uuid.UUID
	UserID       int
	RefreshToken string
	FingerPrint FingerPrint
	ExpireAt     time.Time
	CreatedAt    *time.Time
	ClosedAt     *time.Time
	UpdatedAt    *time.Time
}