package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/domain/models"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/storage/postgres/entities"
	"github.com/jbakhtin/goph-keeper/internal/storage/postgres/query"
)

var _ secondary_ports.SessionRepository = &SessionRepository{}

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

func (s *SessionRepository) Get(ctx context.Context, id int) (*models.Session, error) {
	var entity entities.Session
	err := s.QueryRowContext(ctx, query.GetSessionByID, &id).
		Scan(&entity.ID,
			&entity.UserID,
			&entity.RefreshToken,
			&entity.FingerPrint,
			&entity.ExpireAt,
			&entity.CreatedAt,
			&entity.ClosedAt,
			&entity.UpdatedAt)
	if err != nil {
		return nil, err
	}

	model := entity.ToModel()
	return &model, nil
}

func (s *SessionRepository) Create(ctx context.Context, model models.Session) (*models.Session, error) {
	entity := entities.NewEntity(model)

	err := s.QueryRowContext(ctx, query.CreateSession, entity.UserID, entity.FingerPrint, entity.ExpireAt).
		Scan(&entity.ID,
			&entity.UserID,
			&entity.RefreshToken,
			&entity.FingerPrint,
			&entity.ExpireAt,
			&entity.CreatedAt,
			&entity.ClosedAt,
			&entity.UpdatedAt)
	if err != nil {
		return nil, err
	}

	model = entity.ToModel()

	return &model, nil
}

func (s *SessionRepository) Update(ctx context.Context, session models.Session) (*models.Session, error) {
	var entity entities.Session
	err := s.QueryRowContext(ctx, query.UpdateSessionByID, session.ID, session.UserID, session.FingerPrint, session.ExpireAt, session.ClosedAt).
		Scan(&entity.ID,
			&entity.UserID,
			&entity.RefreshToken,
			&entity.FingerPrint,
			&entity.ExpireAt,
			&entity.CreatedAt,
			&entity.ClosedAt,
			&entity.UpdatedAt)
	if err != nil {
		return nil, err
	}

	model := entity.ToModel()
	return &model, nil
}

func (s *SessionRepository) Delete(ctx context.Context, session models.Session) (*models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionRepository) Search(ctx context.Context, specs queryspec.QuerySpecification) ([]*models.Session, error) {
	rows, err := s.QueryContext(ctx, fmt.Sprintf("%s %s", query.SearchSessionsTemp, specs.Query()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*models.Session, 0)
	for rows.Next() {
		var entity entities.Session
		err = rows.Scan(&entity.ID,
			&entity.UserID,
			&entity.RefreshToken,
			&entity.FingerPrint,
			&entity.ExpireAt,
			&entity.CreatedAt,
			&entity.ClosedAt,
			&entity.UpdatedAt)
		if err != nil {
			return nil, err
		}

		session := entity.ToModel()
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
