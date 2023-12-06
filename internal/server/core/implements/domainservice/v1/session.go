package domainservice

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/domainservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/ports/output/repositories/v1"
)

var _ domainservices.SessionDomainServiceInterface = &sessionDomainService{}

type sessionDomainService struct {
	repo repositories.SessionRepositoryInterface
	cfg  config.Config
}

func NewSessionDomainService(cfg config.Config, repo repositories.SessionRepositoryInterface) (*sessionDomainService, error) {
	return &sessionDomainService{
		cfg:  cfg,
		repo: repo,
	}, nil
}

func (s *sessionDomainService) CreateSession(ctx context.Context, UserId types.Id, FingerPrint types.FingerPrint, ExpireAt types.TimeStamp) (*models.Session, error) {
	session, err := s.repo.SaveSession(ctx, UserId, FingerPrint, ExpireAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionDomainService) GetSessionByID(ctx context.Context, id types.Id) (*models.Session, error) {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionDomainService) GetSessionByRefreshToken(ctx context.Context, refreshToken types.RefreshToken) (*models.Session, error) {
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionDomainService) GetSessionByFingerPrintAndUserId(ctx context.Context, fingerPrint types.FingerPrint, userId types.Id) (*models.Session, error) {
	session, err := s.repo.GetSessionByUserIDAndFingerPrint(ctx, userId, fingerPrint)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionDomainService) CloseSession(ctx context.Context, session models.Session) (*models.Session, error) {
	closed, err := s.repo.CloseSessionByID(ctx, *session.ID)
	if err != nil {
		return nil, err
	}

	return closed, nil
}

func (s *sessionDomainService) UpdateRefreshToken(ctx context.Context, session models.Session) (*models.Session, error) {
	if session.IsClosed() {
		return nil, errors.New("session is closed")
	}

	if session.IsExpired() {
		return nil, errors.New("session is expired")
	}

	closed, err := s.repo.UpdateRefreshTokenById(ctx, *session.ID)
	if err != nil {
		return nil, err
	}

	return closed, nil
}
