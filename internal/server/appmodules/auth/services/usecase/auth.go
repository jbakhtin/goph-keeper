package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
	primary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/primary"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/webserver/grpc/interceptors"

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

type Service struct {
	cfg                       Config
	lgr                       secondary_ports.Logger
	passwordAppService        PasswordService
	accessTokenAppService     AccessTokenService
	sessionRepository         secondary_ports.SessionRepository
	sessionQuerySpecification secondary_ports.SessionQuerySpecification
	userQuerySpecification    secondary_ports.UserQuerySpecification
	userRepository            secondary_ports.UserRepository
}

func New(
	cfg Config,
	lgr secondary_ports.Logger,
	passwordAppService PasswordService,
	accessTokenAppService AccessTokenService,
	sessionRepository secondary_ports.SessionRepository,
	sessionQuerySpecification secondary_ports.SessionQuerySpecification,
	userQuerySpecification secondary_ports.UserQuerySpecification,
	userRepository secondary_ports.UserRepository) (*Service, error) {
	return &Service{
		cfg:                       cfg,
		lgr:                       lgr,
		passwordAppService:        passwordAppService,
		accessTokenAppService:     accessTokenAppService,
		sessionRepository:         sessionRepository,
		sessionQuerySpecification: sessionQuerySpecification,
		userQuerySpecification:    userQuerySpecification,
		userRepository:            userRepository,
	}, nil
}

func (us *Service) RegisterUser(ctx context.Context, email, rawPassword string) (*models.User, error) {
	users, err := us.userRepository.Search(ctx, us.userQuerySpecification.Limit(
		us.userQuerySpecification.Where(
			us.userQuerySpecification.Email(email),
		), 1,
	))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "get user by email")
	}

	if len(users) != 0 {
		return nil, apperror.ErrUserAlreadyExists
	}

	hashedPassword, err := us.passwordAppService.HashPassword(rawPassword)
	if err != nil {
		return nil, errors.Wrap(err, "hash password")
	}

	user := &models.User{
		Email:    email,
		Password: hashedPassword,
	}

	user, err = us.userRepository.Create(ctx, *user) // ToDo: need to pass model
	if err != nil {
		return nil, errors.Wrap(err, "save new user")
	}

	return user, nil
}

func (us *Service) LoginUser(ctx context.Context, email string, password string, fingerPrint models.FingerPrint) (*types.TokensPair, error) {
	users, err := us.userRepository.Search(ctx, us.userQuerySpecification.Limit(
		us.userQuerySpecification.Where(
			us.userQuerySpecification.Email(email),
		), 1,
	))
	if err != nil {
		return nil, errors.Wrap(err, "get user by email")
	}

	if len(users) == 0 {
		return nil, apperror.ErrUserNotFound
	}

	user := users[0]

	if ok, err := us.passwordAppService.CheckPassword(password, user.Password); !ok {
		return nil, errors.New("check password")
	} else if err != nil {
		return nil, errors.Wrap(err, "check password")
	}

	session, err := us.sessionRepository.Create(ctx, models.Session{
		UserID:      *user.ID,
		FingerPrint: fingerPrint,
		ExpireAt:    time.Now().Add(us.cfg.GetSessionExpire()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "create session") // ToDo: mistake
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

func (us *Service) RefreshToken(ctx context.Context, refreshToken string) (*types.TokensPair, error) {
	sessions, err := us.sessionRepository.Search(ctx, us.sessionQuerySpecification.Limit(
		us.sessionQuerySpecification.Where(
			us.sessionQuerySpecification.RefreshToken(refreshToken),
		), 10,
	))
	if err != nil {
		return nil, errors.Wrap(err, "get session by refresh_token")
	}

	if len(sessions) == 0 {
		return nil, errors.Wrap(err, "sessions not found")
	}

	session := sessions[0]

	session.UpdateRefreshToken()
	session, err = us.sessionRepository.Update(ctx, *session)
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

func (us *Service) Logout(ctx context.Context, logOutType primary_ports.LogOutType) (sessions []*models.Session, err error) {
	var sessionID = ctx.Value(interceptors.ContextKeySessionID)
	var userID = ctx.Value(interceptors.ContextKeyUserID)

	// ToDo: добавить проверку на истечение срока жизни сессии и то что сессия уже закрыта
	switch logOutType {
	case primary_ports.LogoutTypeThis:
		session, err := us.sessionRepository.Get(ctx, sessionID.(int))
		if err != nil {
			return nil, errors.Wrap(err, "get session by session id")
		}

		session.Close()
		session, err = us.sessionRepository.Update(ctx, *session)
		if err != nil {
			return nil, errors.Wrap(err, "close current session by session_id")
		}
		sessions = append(sessions, session)
	case primary_ports.LogoutTypeAll:
		sessions, err = us.sessionRepository.Search(ctx, us.sessionQuerySpecification.Where(
			us.sessionQuerySpecification.And(
				us.sessionQuerySpecification.UserID(userID.(int)),
				us.sessionQuerySpecification.IsNotClosed(),
			),
		))
		if err != nil {
			return nil, errors.Wrap(err, "get sessions by user_id")
		}

		for index, session := range sessions {
			session.Close()
			session, err = us.sessionRepository.Update(ctx, *session)
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
