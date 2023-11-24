package session

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/application/types"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/repositories"
	"github.com/jbakhtin/goph-keeper/internal/server/infastructure/database/postgres"
	"github.com/jbakhtin/goph-keeper/internal/server/infastructure/database/postgres/query"
)

var _ repositories.SessionRepository = (*Repository)(nil)

type Repository struct {
	*postgres.Postgres
}

func New(client postgres.Postgres) (*Repository, error) {
	return &Repository{
		&client,
	}, nil
}

func (r *Repository) SaveSession(ctx context.Context, session models.Session) (*models.Session, error) {
	var created models.Session
	err := r.QueryRowContext(ctx, query.CreateSession, &session.UserId, &session.FingerPrint, &session.ExpireAt).
		Scan(&created.ID,
			&created.UserId,
			&created.RefreshToken,
		created.FingerPrint,
		&created.ExpireAt,
		&created.CreatedAt,
		&created.ClosedAt,
		&created.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *Repository) UpdateSessionByID(ctx context.Context, id types.Id, session models.Session) (*models.Session, error) {
	var updated models.Session
	err := r.QueryRowContext(ctx, query.UpdateSessionByID, id, session.ExpireAt).
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

func (r *Repository) GetSessionByUserIDAndFingerPrint(ctx context.Context, userId types.Id, fingerPrint types.FingerPrint) (*models.Session, error) {
	var got models.Session
	err := r.QueryRowContext(ctx, query.GetSessionByUserIDAndFingerPrint, userId, fingerPrint).
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


func (r *Repository) UpdateRefreshTokenById(ctx context.Context, id types.Id) (*models.Session, error) {
	var withNewRefresh models.Session
	err := r.QueryRowContext(ctx, query.UpdateSessionRefreshTokenById, id).
		Scan(&withNewRefresh.ID,
			&withNewRefresh.UserId,
			&withNewRefresh.RefreshToken,
			withNewRefresh.FingerPrint,
			&withNewRefresh.ExpireAt,
			&withNewRefresh.CreatedAt,
			&withNewRefresh.ClosedAt,
			&withNewRefresh.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &withNewRefresh, nil
}

func (r *Repository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var get models.Session
	err := r.QueryRowContext(ctx, query.GetSessionByRefreshToken, &refreshToken).
		Scan(&get.ID,
			&get.UserId,
			&get.RefreshToken,
			get.FingerPrint,
			&get.ExpireAt,
			&get.CreatedAt,
			&get.ClosedAt,
			&get.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &get, nil
}


func (r *Repository) CloseSessionByID(ctx context.Context, id types.Id) (*models.Session, error) {
	var closed models.Session
	err := r.QueryRowContext(ctx, query.CloseSessionByID, id).
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

func (r *Repository) CloseSessionsByUserId(ctx context.Context, userId types.Id) ([]*models.Session, error) {
	rows, err := r.QueryContext(ctx, query.CloseSessionsByUserID, userId)
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