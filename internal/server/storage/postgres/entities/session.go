package entities

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"time"
)

var _ sql.Scanner = &FingerPrint{}

// MyScanner - пользовательский тип для реализации интерфейса sql.Scanner
type FingerPrint map[string]any

// Scan - метод, реализующий интерфейс sql.Scanner
func (m *FingerPrint) Scan(value interface{}) error {
	// Проверяем, что значение является срезом байтов ([]uint8)
	if value == nil {
		*m = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, received %T", value)
	}

	// Декодируем срез байтов в map[string]interface{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	return nil
}

type Session struct {
	ID           *int        `json:"id,omitempty"`
	UserID       int         `json:"user_id,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	FingerPrint  FingerPrint `json:"finger_print,omitempty"`
	ExpireAt     time.Time   `json:"expire_at"`
	CreatedAt    *time.Time  `json:"created_at,omitempty"`
	ClosedAt     *time.Time  `json:"closed_at,omitempty"`
	UpdatedAt    *time.Time  `json:"updated_at,omitempty"`
}

func NewEntity(model models.Session) Session {
	return Session{
		ID:           model.ID,
		UserID:       model.UserID,
		RefreshToken: model.RefreshToken,
		FingerPrint:  FingerPrint(model.FingerPrint),
		ExpireAt:     model.ExpireAt,
		CreatedAt:    model.CreatedAt,
		ClosedAt:     model.ClosedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func (s *Session) ToModel() models.Session {
	return models.Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		FingerPrint:  models.FingerPrint(s.FingerPrint),
		ExpireAt:     s.ExpireAt,
		CreatedAt:    s.CreatedAt,
		ClosedAt:     s.ClosedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
