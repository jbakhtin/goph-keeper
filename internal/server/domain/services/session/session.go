package session

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jbakhtin/goph-keeper/internal/server/application/types"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/repositories"
	"github.com/pkg/errors"
	"time"
)

type TokensPair struct {
	AccessToken string
	RefreshToken string
}

type IConfig interface {
	GetAppKey() string
}

type Service struct {
	cfg    IConfig
	repo   repositories.SessionRepository
}

func New(cfg IConfig, repo repositories.SessionRepository) (*Service, error) {
	return &Service{
		cfg: cfg,
		repo:   repo,
	}, nil
}

func (s *Service) Create(ctx context.Context, user models.User, fingerPrint types.FingerPrint) (*TokensPair, error) {
	session, err := s.repo.GetSessionByUserIDAndFingerPrint(ctx, *user.ID, fingerPrint)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "get session")
		}
	} else {
		_, err = s.repo.CloseSessionByID(ctx, *session.ID)
		if err != nil {
			return nil, errors.Wrap(err, "close session")
		}
	}

	expireAt := time.Now().Add(time.Hour * 24 * 30) // ToDo: move expire parameters to config
	session, err = s.repo.SaveSession(ctx, models.Session{
		UserId: *user.ID,
		ExpireAt: expireAt,
		FingerPrint: &fingerPrint,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create session")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": jwt.MapClaims{
			"user_id": user.ID,
			"session_id": session.ID,
		},
		"exp": time.Minute * 10,
	})

	accessToken, err := token.SignedString([]byte(s.cfg.GetAppKey()))
	if err != nil {
		return nil, errors.Wrap(err, "sign auth jwt")
	}

	return &TokensPair{
		AccessToken: accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (s *Service) Update(ctx context.Context, refreshToken string) (*TokensPair, error) {
	fmt.Println(refreshToken)
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByRefreshToken")
	}

	session, err = s.repo.UpdateRefreshTokenById(ctx, *session.ID)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateRefreshTokenById")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": jwt.MapClaims{
			"user_id": session.UserId,
			"session_id": session.ID,
		},
		"exp": time.Minute * 10,
	})

	accessToken, err := token.SignedString([]byte(s.cfg.GetAppKey()))
	if err != nil {
		return nil, errors.Wrap(err, "sign auth jwt")
	}

	return &TokensPair{
		AccessToken: accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}


func (s *Service) Close(ctx context.Context, logoutType models.LogoutType) (sessions []*models.Session, err error) {
	switch logoutType {
	case models.LogoutType_ALL:
		if userId, ok := ctx.Value(types.ContextKeyUserID).(types.Id); ok {
			sessions, err = s.repo.CloseSessionsByUserId(ctx, userId)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("close sessions by user_id: %v", userId))
			}
		} else {
			return nil, errors.New("session not available")
		}
	case models.LogoutType_THIS:
		if sessionId, ok := ctx.Value(types.ContextKeySessionID).(types.Id); ok {
			var session *models.Session
			session, err = s.repo.CloseSessionByID(ctx, sessionId)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("close session by session_id: %v", sessionId))
			}
			sessions = append(sessions, session)
		} else {
			return nil, errors.New("session not available")
		}
	}

	return sessions, nil
}