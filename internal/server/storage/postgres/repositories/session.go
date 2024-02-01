package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/query"
	"strings"
	"time"
)

var _ secondary_ports.SessionRepository = &SessionRepository{}
var _ secondary_ports.SessionSpecifications = &SessionSpecifications{}

// Specifications

// Where Specification

type WhereSpecification struct {
	specification secondary_ports.Specification
}

func (s *WhereSpecification) Query() string {
	return fmt.Sprintf("WHERE %s", s.specification.Query())
}

func (s *WhereSpecification) Value() []any {
	return s.specification.Value()
}

// UserID Specification

type LimitSpecification struct {
	limit         int
	specification secondary_ports.Specification
}

func (s *LimitSpecification) Query() string {
	return fmt.Sprintf("%s LIMIT %v", s.specification.Query(), s.limit)
}

func (s *LimitSpecification) Value() []any {
	return append(s.specification.Value(), s.limit)
}

// UserID Specification

type UserIDSpecification struct {
	userID int
}

func (s *UserIDSpecification) Query() string {
	return fmt.Sprintf("user_id = %v", s.userID)
}

func (s *UserIDSpecification) Value() []any {
	return []any{s.userID}
}

// FingerPrint Specification

type FingerPrintSpecification struct {
	fingerPrint string
}

func (s *FingerPrintSpecification) Query() string {
	return fmt.Sprintf("finger_print = %v", s.fingerPrint)
}

func (s *FingerPrintSpecification) Value() []any {
	return []any{s.fingerPrint}
}

// RefreshToken Specification

type RefreshTokenSpecification struct {
	refreshToken string
}

func (s *RefreshTokenSpecification) Query() string {
	return fmt.Sprintf("refresh_token = %v", s.refreshToken)
}

func (s *RefreshTokenSpecification) Value() []any {
	return []any{s.refreshToken}
}

// IsNotClosed Specification

type IsNotClosedSpecification struct{}

func (s *IsNotClosedSpecification) Query() string {
	return "closed_at IS NULL"
}

func (s *IsNotClosedSpecification) Value() []any {
	return nil
}

// Or Specification

type OrSpecification struct {
	specifications []secondary_ports.Specification
}

func (s *OrSpecification) Query() string {
	var queries []string
	for _, specification := range s.specifications {
		queries = append(queries, specification.Query())
	}

	query := strings.Join(queries, " OR ")

	return fmt.Sprintf("(%s)", query)
}

func (s *OrSpecification) Value() []any {
	var values []interface{}
	for _, specification := range s.specifications {
		values = append(values, specification.Value()...)
	}
	return values
}

// And Specification

type AndSpecification struct {
	specifications []secondary_ports.Specification
}

func (s *AndSpecification) Query() string {
	var queries []string
	for _, specification := range s.specifications {
		queries = append(queries, specification.Query())
	}

	query := strings.Join(queries, " AND ")

	return fmt.Sprintf("%s", query)
}

func (s *AndSpecification) Value() []any {
	var values []interface{}
	for _, specification := range s.specifications {
		values = append(values, specification.Value()...)
	}
	return values
}

// Session Specifications

type SessionSpecifications struct {
	specifications []secondary_ports.Specification
}

func NewSessionSpecifications() (SessionSpecifications, error) {
	return SessionSpecifications{
		specifications: make([]secondary_ports.Specification, 0),
	}, nil
}

func (s SessionSpecifications) Limit(specification secondary_ports.Specification, i int) secondary_ports.Specification {
	return &LimitSpecification{
		specification: specification,
		limit:         i,
	}
}

func (s SessionSpecifications) Where(specification secondary_ports.Specification) secondary_ports.Specification {
	return &WhereSpecification{
		specification: specification,
	}
}

func (s SessionSpecifications) Or(specifications ...secondary_ports.Specification) secondary_ports.Specification {
	return &OrSpecification{
		specifications: specifications,
	}
}

func (s SessionSpecifications) And(specifications ...secondary_ports.Specification) secondary_ports.Specification {
	return &AndSpecification{
		specifications: specifications,
	}
}

func (s SessionSpecifications) UserID(userId int) secondary_ports.Specification {
	return &UserIDSpecification{
		userID: userId,
	}
}

func (s SessionSpecifications) IsNotClosed() secondary_ports.Specification {
	return &IsNotClosedSpecification{}
}

func (s SessionSpecifications) FingerPrint(fingerPrint models.FingerPrint) secondary_ports.Specification {
	buf, _ := json.Marshal(fingerPrint)
	return &FingerPrintSpecification{
		fingerPrint: string(buf),
	}
}

func (s SessionSpecifications) RefreshToken(refreshToken string) secondary_ports.Specification {
	return &RefreshTokenSpecification{
		refreshToken: refreshToken,
	}
}

// Repository

type SessionRepository struct {
	*sql.DB
	lgr *zap.Logger
}

func NewSessionRepository(lgr *zap.Logger, client *sql.DB) (*SessionRepository, error) {
	return &SessionRepository{
		DB:  client,
		lgr: lgr,
	}, nil
}

func (s *SessionRepository) Search(ctx context.Context, specs secondary_ports.Specification) ([]*models.Session, error) {
	query := fmt.Sprintf("%s %s", query.SearchSessionsTemp, specs.Query())

	rows, err := s.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := make([]*models.Session, 0)
	for rows.Next() {
		var session models.Session
		err = rows.Scan(&session.ID,
			&session.UserID,
			&session.RefreshToken,
			&session.FingerPrint,
			&session.ExpireAt,
			&session.CreatedAt,
			&session.ClosedAt,
			&session.UpdatedAt)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *SessionRepository) SaveSession(ctx context.Context, UserID int, FingerPrint models.FingerPrint, ExpireAt time.Time) (*models.Session, error) {
	var created models.Session
	err := s.QueryRowContext(ctx, query.CreateSession, UserID, FingerPrint, ExpireAt).
		Scan(&created.ID,
			&created.UserID,
			&created.RefreshToken,
			&created.FingerPrint,
			&created.ExpireAt,
			&created.CreatedAt,
			&created.ClosedAt,
			&created.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *SessionRepository) UpdateSession(ctx context.Context, session models.Session) (*models.Session, error) {
	var updated models.Session
	err := s.QueryRowContext(ctx, query.UpdateSessionByID, session.ID, session.ExpireAt).
		Scan(&updated.ID,
			&updated.UserID,
			&updated.RefreshToken,
			&updated.FingerPrint,
			&updated.ExpireAt,
			&updated.CreatedAt,
			&updated.ClosedAt,
			&updated.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (s *SessionRepository) GetSessionByUserIDAndFingerPrint(ctx context.Context, userID int, fingerPrint models.FingerPrint) (*models.Session, error) {
	var got models.Session
	err := s.QueryRowContext(ctx, query.GetSessionByUserIDAndFingerPrint, userID, fingerPrint).
		Scan(&got.ID,
			&got.UserID,
			&got.RefreshToken,
			got.FingerPrint,
			&got.ExpireAt,
			&got.CreatedAt,
			&got.ClosedAt,
			&got.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &got, nil
}

func (s *SessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var get models.Session
	err := s.QueryRowContext(ctx, query.GetSessionByRefreshToken, &refreshToken).
		Scan(&get.ID,
			&get.UserID,
			&get.RefreshToken,
			&get.FingerPrint,
			&get.ExpireAt,
			&get.CreatedAt,
			&get.ClosedAt,
			&get.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &get, nil
}

func (s *SessionRepository) GetSession(ctx context.Context, id int) (*models.Session, error) {
	var get models.Session
	err := s.QueryRowContext(ctx, query.GetSessionByID, &id).
		Scan(&get.ID,
			&get.UserID,
			&get.RefreshToken,
			&get.FingerPrint,
			&get.ExpireAt,
			&get.CreatedAt,
			&get.ClosedAt,
			&get.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &get, nil
}
