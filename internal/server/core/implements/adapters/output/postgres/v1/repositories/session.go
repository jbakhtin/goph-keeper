package repositories

import (
	"context"
	"fmt"

	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/output/postgres/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/output/postgres/v1/query"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/ports/output/repositories/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

var _ repositories.SessionRepositoryInterface = (*SessionRepository)(nil)

type SessionRepository struct {
	*postgres.Postgres
}

func NewSessionRepository(client postgres.Postgres) (*SessionRepository, error) { // ToDo: need remove client parameter
	return &SessionRepository{
		&client,
	}, nil
}

func (s *SessionRepository) SaveSession(ctx context.Context, UserId types.Id, FingerPrint types.FingerPrint, ExpireAt types.TimeStamp) (*models.Session, error) {
	var created models.Session
	err := s.QueryRowContext(ctx, query.CreateSession, UserId, FingerPrint, ExpireAt).
		Scan(&created.ID,
			&created.UserId,
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

func (s *SessionRepository) UpdateSessionByID(ctx context.Context, id types.Id, session models.Session) (*models.Session, error) {
	var updated models.Session
	err := s.QueryRowContext(ctx, query.UpdateSessionByID, id, session.ExpireAt).
		Scan(&updated.ID,
			&updated.UserId,
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

func (s *SessionRepository) GetSessionsByUserID(ctx context.Context, userId types.Id) ([]*models.Session, error) {
	fmt.Println(userId)
	rows, err := s.QueryContext(ctx, query.GetSessionByUserID, userId)
	if err != nil {
		return nil, err
	}

	sessions := make([]*models.Session, 0)
	for rows.Next() {
		var session models.Session
		err = rows.Scan(&session.ID,
			&session.UserId,
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

	return sessions, nil
}

func (s *SessionRepository) GetSessionByUserIDAndFingerPrint(ctx context.Context, userId types.Id, fingerPrint types.FingerPrint) (*models.Session, error) {
	var got models.Session
	err := s.QueryRowContext(ctx, query.GetSessionByUserIDAndFingerPrint, userId, fingerPrint).
		Scan(&got.ID,
			&got.UserId,
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

func (s *SessionRepository) UpdateRefreshTokenById(ctx context.Context, id types.Id) (*models.Session, error) {
	var withNewRefresh models.Session
	err := s.QueryRowContext(ctx, query.UpdateSessionRefreshTokenById, id).
		Scan(&withNewRefresh.ID,
			&withNewRefresh.UserId,
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

func (s *SessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken types.RefreshToken) (*models.Session, error) {
	var get models.Session
	err := s.QueryRowContext(ctx, query.GetSessionByRefreshToken, &refreshToken).
		Scan(&get.ID,
			&get.UserId,
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

func (s *SessionRepository) GetSessionByID(ctx context.Context, id types.Id) (*models.Session, error) {
	var get models.Session
	err := s.QueryRowContext(ctx, query.GetSessionByID, &id).
		Scan(&get.ID,
			&get.UserId,
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

func (s *SessionRepository) CloseSessionByID(ctx context.Context, id types.Id) (*models.Session, error) {
	var closed models.Session
	err := s.QueryRowContext(ctx, query.CloseSessionByID, id).
		Scan(&closed.ID,
			&closed.UserId,
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

func (s *SessionRepository) CloseSessionsByUserId(ctx context.Context, userId types.Id) ([]*models.Session, error) {
	rows, err := s.QueryContext(ctx, query.CloseSessionsByUserID, userId)
	if err != nil {
		return nil, err
	}

	var sessions []*models.Session
	for rows.Next() {
		var session models.Session
		err = rows.Scan(&session.ID,
			&session.UserId,
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

	return sessions, nil
}
