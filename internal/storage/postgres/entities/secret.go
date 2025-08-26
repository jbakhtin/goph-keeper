package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/go-faster/errors"
)

type Secret struct {
	ID          int
	UserID      int
	Description string
	Data        Data
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Data map[string]any

func (a Data) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Data) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
