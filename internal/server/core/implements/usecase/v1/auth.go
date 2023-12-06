package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/ports/output/repositories/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/apperror"
	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/appservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/domainservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/usecases/v1"
	"github.com/pkg/errors"
)

var _ usecases.AuthUseCaseInterface = &AuthUseCase{}

type AuthUseCase struct {
	userDomainService     domainservices.UserDomainServiceInterface
	sessionDomainService  domainservices.SessionDomainServiceInterface
	passwordAppService    appservices.PasswordAppServiceInterface
	accessTokenAppService appservices.AccessTokenAppServiceInterface
	sessionsRepository    repositories.SessionRepositoryInterface
	cfg                   config.Config
}

func NewAuthUseCase(
	cfg config.Config,
	userDomainService domainservices.UserDomainServiceInterface,
	sessionDomainService domainservices.SessionDomainServiceInterface,
	passwordAppService appservices.PasswordAppServiceInterface,
	accessTokenAppService appservices.AccessTokenAppServiceInterface,
	sessionsRepository repositories.SessionRepositoryInterface) (*AuthUseCase, error) {
	return &AuthUseCase{
		userDomainService:     userDomainService,
		passwordAppService:    passwordAppService,
		accessTokenAppService: accessTokenAppService,
		sessionDomainService:  sessionDomainService,
		sessionsRepository:    sessionsRepository,
		cfg:                   cfg,
	}, nil
}

func (us *AuthUseCase) RegisterUser(ctx context.Context, email, rawPassword string) (*models.User, error) {
	user, err := us.userDomainService.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "get user by email")
	}

	if user != nil {
		return nil, apperror.UserAlreadyExists
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
	user, err := us.userDomainService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "get user by email")
	}

	if ok, err := us.passwordAppService.CheckPassword(password, user.Password); ok == false {
		return nil, errors.New("check password")
	} else if err != nil {
		return nil, errors.Wrap(err, "check password")
	}

	session, err := us.sessionDomainService.CreateSession(ctx, *user.ID, fingerPrint, types.TimeStamp(time.Now().Add(time.Hour*24*30))) // ToDo: move expired_at to cfg
	if err != nil {
		return nil, errors.Wrap(err, "create session") //ToDo: mistake
	}

	accessToken, err := us.accessTokenAppService.Create(*user.ID, *session.ID, time.Minute*5) // ToDo: move access_token duration to cfg
	if err != nil {
		return nil, errors.Wrap(err, "create access_token")
	}

	return &types.TokensPair{
		AccessToken:  *accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (us *AuthUseCase) RefreshToken(ctx context.Context, refreshToken types.RefreshToken) (*types.TokensPair, error) {
	session, err := us.sessionDomainService.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "get session by refresh_token")
	}

	session, err = us.sessionDomainService.UpdateRefreshToken(ctx, *session)
	if err != nil {
		return nil, errors.Wrap(err, "update refresh token")
	}

	accessToken, err := us.accessTokenAppService.Create(session.UserId, *session.ID, time.Minute*5) // ToDo: mode access_token duration to cfg
	if err != nil {
		return nil, errors.Wrap(err, "create access_token")
	}

	return &types.TokensPair{
		AccessToken:  *accessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (us *AuthUseCase) Logout(ctx context.Context, logoutType types.LogoutType) (sessions []*models.Session, err error) {
	var sessionId = ctx.Value(types.ContextKeySessionID)
	var userId = ctx.Value(types.ContextKeyUserID)

	fmt.Println("LogoutType_TYPE_ALL")

	switch logoutType {
	case types.LogoutType_THIS:
		session, err := us.sessionDomainService.GetSessionByID(ctx, sessionId.(types.Id))
		if err != nil {
			return nil, errors.Wrap(err, "get session by session id")
		}

		session, err = us.sessionDomainService.CloseSession(ctx, *session)
		if err != nil {
			return nil, errors.Wrap(err, "close current session by session_id")
		}
		sessions = append(sessions, session)
	case types.LogoutType_ALL:
		sessions, err = us.sessionsRepository.GetSessionsByUserID(ctx, userId.(types.Id))

		for index, session := range sessions {
			session, err = us.sessionDomainService.CloseSession(ctx, *session)
			sessions[index] = session
			if err != nil {
				return nil, errors.Wrap(err, "close all sessions by user_id")
			}
		}
	}

	return sessions, nil
}
