package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/query"
	"time"
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

func (s *SessionRepository) UpdateSessionByID(ctx context.Context, id int, session models.Session) (*models.Session, error) {
	var updated models.Session
	err := s.QueryRowContext(ctx, query.UpdateSessionByID, id, session.ExpireAt).
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

func (s *SessionRepository) GetSessionsByUserID(ctx context.Context, userID int) ([]*models.Session, error) {
	fmt.Println(userID)
	rows, err := s.QueryContext(ctx, query.GetSessionByUserID, userID)
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

func (s *SessionRepository) UpdateRefreshTokenByID(ctx context.Context, id int) (*models.Session, error) {
	var withNewRefresh models.Session
	err := s.QueryRowContext(ctx, query.UpdateSessionRefreshTokenByID, id).
		Scan(&withNewRefresh.ID,
			&withNewRefresh.UserID,
			&withNewRefresh.RefreshToken,
			&withNewRefresh.FingerPrint,
			&withNewRefresh.ExpireAt,
			&withNewRefresh.CreatedAt,
			&withNewRefresh.ClosedAt,
			&withNewRefresh.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &withNewRefresh, nil
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

func (s *SessionRepository) GetSessionByID(ctx context.Context, id int) (*models.Session, error) {
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

func (s *SessionRepository) CloseSessionByID(ctx context.Context, id int) (*models.Session, error) {
	var closed models.Session
	err := s.QueryRowContext(ctx, query.CloseSessionByID, id).
		Scan(&closed.ID,
			&closed.UserID,
			&closed.RefreshToken,
			&closed.FingerPrint,
			&closed.ExpireAt,
			&closed.CreatedAt,
			&closed.ClosedAt,
			&closed.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &closed, nil
}

func (s *SessionRepository) CloseSessionsByUserID(ctx context.Context, userID int) ([]*models.Session, error) {
	rows, err := s.QueryContext(ctx, query.CloseSessionsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
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
