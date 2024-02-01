package usecase

import (
	"context"
	"database/sql"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/primary"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/webserver/grpc/interceptors"
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/apperror"
	"github.com/pkg/errors"
)

type Config interface {
	GetAccessTokenExpire() time.Duration
	GetSessionExpire() time.Duration
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, need string) (bool, error)
}

type AccessTokenService interface {
	Create(userID int, sessionID int, duration time.Duration) (*string, error)
}

type AuthUseCase struct {
	cfg                   Config
	lgr                   secondary_ports.Logger
	passwordAppService    PasswordService
	accessTokenAppService AccessTokenService
	sessionRepository     secondary_ports.SessionRepository
	sessionSpecifications secondary_ports.SessionSpecifications
	userRepository        secondary_ports.UserRepository
}

func NewAuthUseCase(
	cfg Config,
	lgr secondary_ports.Logger,
	passwordAppService PasswordService,
	accessTokenAppService AccessTokenService,
	sessionRepository secondary_ports.SessionRepository,
	userRepository secondary_ports.UserRepository) (*AuthUseCase, error) {
	return &AuthUseCase{
		cfg:                   cfg,
		lgr:                   lgr,
		passwordAppService:    passwordAppService,
		accessTokenAppService: accessTokenAppService,
		sessionRepository:     sessionRepository,
		userRepository:        userRepository,
	}, nil
}

func (us *AuthUseCase) RegisterUser(ctx context.Context, email, rawPassword string) (*models.User, error) {
	user, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "get user by email")
	}

	if user != nil {
		return nil, apperror.ErrUserAlreadyExists
	}

	hashedPassword, err := us.passwordAppService.HashPassword(rawPassword)
	if err != nil {
		return nil, errors.Wrap(err, "hash password")
	}

	user = &models.User{
		Email:    email,
		Password: hashedPassword,
	}

	user, err = us.userRepository.SaveUser(ctx, *user) // ToDo: need to pass model
	if err != nil {
		return nil, errors.Wrap(err, "save new user")
	}

	return user, nil
}

func (us *AuthUseCase) LoginUser(ctx context.Context, email string, password string, fingerPrint models.FingerPrint) (*types.TokensPair, error) {
	user, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "get user by email")
	}

	if ok, err := us.passwordAppService.CheckPassword(password, user.Password); !ok {
		return nil, errors.New("check password")
	} else if err != nil {
		return nil, errors.Wrap(err, "check password")
	}

	session, err := us.sessionRepository.SaveSession(ctx, *user.ID, fingerPrint, time.Now().Add(us.cfg.GetSessionExpire()))
	if err != nil {
		return nil, errors.Wrap(err, "create session") //ToDo: mistake
	}

	accessToken, err := us.accessTokenAppService.Create(*user.ID, *session.ID, us.cfg.GetAccessTokenExpire())
	if err != nil {
		return nil, errors.Wrap(err, "create access_token")
	}

	return &types.TokensPair{
		AccessToken:  *accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (us *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*types.TokensPair, error) {
	session, err := us.sessionRepository.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "get session by refresh_token")
	}

	session.UpdateRefreshToken()
	session, err = us.sessionRepository.UpdateSession(ctx, *session)
	if err != nil {
		return nil, errors.Wrap(err, "update refresh token")
	}

	accessToken, err := us.accessTokenAppService.Create(session.UserID, *session.ID, us.cfg.GetAccessTokenExpire()) // ToDo: mode access_token duration to cfg
	if err != nil {
		return nil, errors.Wrap(err, "create access_token")
	}

	return &types.TokensPair{
		AccessToken:  *accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (us *AuthUseCase) Logout(ctx context.Context, logOutType primary_ports.LogOutType) (sessions []*models.Session, err error) {
	var sessionID = ctx.Value(interceptors.ContextKeySessionID)
	var userID = ctx.Value(interceptors.ContextKeyUserID)

	// ToDo: добавить проверку на истечение срока жизни сессии и то что сессия уже закрыта
	switch logOutType {
	case primary_ports.LogoutTypeThis:
		session, err := us.sessionRepository.GetSession(ctx, sessionID.(int))
		if err != nil {
			return nil, errors.Wrap(err, "get session by session id")
		}

		session.Close()
		session, err = us.sessionRepository.UpdateSession(ctx, *session)
		if err != nil {
			return nil, errors.Wrap(err, "close current session by session_id")
		}
		sessions = append(sessions, session)
	case primary_ports.LogoutTypeAll:
		sessions, err = us.sessionRepository.Search(ctx, us.sessionSpecifications.Where(
			us.sessionSpecifications.And(
				us.sessionSpecifications.UserID(userID.(int)),
				us.sessionSpecifications.IsNotClosed(),
			),
		))
		if err != nil {
			return nil, errors.Wrap(err, "get sessions by user_id")
		}

		for index, session := range sessions {
			session.Close()
			session, err = us.sessionRepository.UpdateSession(ctx, *session)
			if err != nil {
				return nil, errors.Wrap(err, "close all sessions by user_id")
			}

			sessions[index] = session
		}
	default:
		return nil, errors.New("logout type error")
	}

	return sessions, nil
}
