package domainservice

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/domainservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/repositories/v1"

	"github.com/go-faster/errors"
)

var _ domainservices.SessionDomainServiceInterface = &sessionDomainService{}

type sessionDomainService struct {
	cfg  config.Interface
	repo repositories.SessionRepositoryInterface
}

func NewSessionDomainService(cfg config.Interface, repo repositories.SessionRepositoryInterface) (*sessionDomainService, error) {
	return &sessionDomainService{
		cfg:  cfg,
		repo: repo,
	}, nil
}

func (s *sessionDomainService) CreateSession(ctx context.Context, UserID types.ID, FingerPrint types.FingerPrint, ExpireAt types.TimeStamp) (*models.Session, error) {
	session, err := s.repo.SaveSession(ctx, UserID, FingerPrint, ExpireAt)
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

	closed, err := s.repo.UpdateRefreshTokenByID(ctx, *session.ID)
	if err != nil {
		return nil, err
	}

	return closed, nil
}
