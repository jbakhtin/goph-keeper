package usecase

import (
	"context"
	"database/sql"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/appservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/domainservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/repositories/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/usecases/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/apperror"
	"github.com/pkg/errors"
)

var _ usecases.AuthUseCaseInterface = &AuthUseCase{}

type AuthUseCase struct {
	cfg                   config.Interface
	lgr logger.Interface
	userDomainService     domainservices.UserDomainServiceInterface
	sessionDomainService  domainservices.SessionDomainServiceInterface
	passwordAppService    appservices.PasswordAppServiceInterface
	accessTokenAppService appservices.AccessTokenAppServiceInterface
	sessionRepository     repositories.SessionRepositoryInterface
	userRepository        repositories.UserRepositoryInterface
}

func NewAuthUseCase(
	cfg config.Interface,
	lgr logger.Interface,
	userDomainService domainservices.UserDomainServiceInterface,
	sessionDomainService domainservices.SessionDomainServiceInterface,
	passwordAppService appservices.PasswordAppServiceInterface,
	accessTokenAppService appservices.AccessTokenAppServiceInterface,
	sessionRepository repositories.SessionRepositoryInterface,
	userRepository repositories.UserRepositoryInterface) (*AuthUseCase, error) {
	return &AuthUseCase{
		cfg:                   cfg,
		lgr: lgr,
		userDomainService:     userDomainService,
		passwordAppService:    passwordAppService,
		accessTokenAppService: accessTokenAppService,
		sessionDomainService:  sessionDomainService,
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

	user, err = us.userDomainService.CreateUser(ctx, email, hashedPassword)
	if err != nil {
		return nil, errors.Wrap(err, "save new user")
	}

	return user, nil
}

func (us *AuthUseCase) LoginUser(ctx context.Context, email string, password string, fingerPrint types.FingerPrint) (*types.TokensPair, error) {
	user, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "get user by email")
	}

	if ok, err := us.passwordAppService.CheckPassword(password, user.Password); !ok {
		return nil, errors.New("check password")
	} else if err != nil {
		return nil, errors.Wrap(err, "check password")
	}

	session, err := us.sessionDomainService.CreateSession(ctx, *user.ID, fingerPrint, types.TimeStamp(time.Now().Add(us.cfg.GetSessionExpire())))
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

func (us *AuthUseCase) RefreshToken(ctx context.Context, refreshToken types.RefreshToken) (*types.TokensPair, error) {
	session, err := us.sessionRepository.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "get session by refresh_token")
	}

	session, err = us.sessionDomainService.UpdateRefreshToken(ctx, *session)
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

func (us *AuthUseCase) Logout(ctx context.Context, logoutType types.LogoutType) (sessions []*models.Session, err error) {
	var sessionID = ctx.Value(types.ContextKeySessionID)
	var userID = ctx.Value(types.ContextKeyUserID)

	switch logoutType {
	case types.LogoutTypeThis:
		session, err := us.sessionRepository.GetSessionByID(ctx, sessionID.(types.ID))
		if err != nil {
			return nil, errors.Wrap(err, "get session by session id")
		}

		session, err = us.sessionDomainService.CloseSession(ctx, *session)
		if err != nil {
			return nil, errors.Wrap(err, "close current session by session_id")
		}
		sessions = append(sessions, session)
	case types.LogoutTypeAll:
		sessions, err = us.sessionRepository.GetSessionsByUserID(ctx, userID.(types.ID))
		if err != nil {
			return nil, errors.Wrap(err, "get sessions by user_id")
		}

		for index, session := range sessions {
			session, err = us.sessionDomainService.CloseSession(ctx, *session)
			if err != nil {
				return nil, errors.Wrap(err, "close all sessions by user_id")
			}

			sessions[index] = session
		}
	}

	return sessions, nil
}
