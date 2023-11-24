package auth

import (
	"github.com/jbakhtin/goph-keeper/internal/server/domain/services/auth"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/services/session"
)

type Builder struct {
	handler *Handler
	authService auth.Service
	sessionService session.Service
	err error
}

func New(authService auth.Service, sessionService session.Service) *Builder {
	return &Builder{
		handler: nil,
		authService: authService,
		sessionService: sessionService,
	}
}

func (b *Builder) Build() (*Handler, error) {
	if b.err != nil {
		return nil, b.err
	}

	b.handler = &Handler{
		authService: b.authService,
		sessionService: b.sessionService,
	}

	return b.handler, nil
}