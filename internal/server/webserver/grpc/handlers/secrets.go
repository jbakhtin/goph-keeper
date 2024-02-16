package handlers

import (
	"context"
	"encoding/json"

	"github.com/jbakhtin/goph-keeper/gen/go/v1/secrets"

	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/domain/models"
	primary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/ports/primary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/webserver/grpc/interceptors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-faster/errors"
)

var _ secrets.SecretsServiceServer = &SecretsHandler{}

type SecretsHandler struct {
	lgr           *zap.Logger
	secretUseCase primary_ports.UseCase
	validator     *protovalidate.Validator

	secrets.UnimplementedSecretsServiceServer
}

func NewKeyValueHandler(lgr *zap.Logger, secretUseCase primary_ports.UseCase) (*SecretsHandler, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &SecretsHandler{
		lgr:           lgr,
		secretUseCase: secretUseCase,
		validator:     validator,
	}, nil
}

func toModel(ctx context.Context, pb *secrets.CrateRequest) (models.Secret, error) {
	var data models.Data
	err := json.Unmarshal([]byte(pb.Data), &data)
	if err != nil {
		return models.Secret{}, err
	}

	return models.Secret{
		UserID:      ctx.Value(interceptors.ContextKeyUserID).(int),
		Type:        pb.Type,
		Data:        data,
		Description: pb.Description,
	}, nil
}

func (h *SecretsHandler) Create(ctx context.Context, request *secrets.CrateRequest) (*emptypb.Empty, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	model, err := toModel(ctx, request)
	if err != nil {
		return nil, err
	}

	_, err = h.secretUseCase.Create(ctx, model)
	if err != nil {
		return nil, errors.Wrap(err, "create key value")
	}

	return &emptypb.Empty{}, nil
}
