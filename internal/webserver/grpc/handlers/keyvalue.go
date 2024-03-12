package handlers

import (
	"context"

	"github.com/jbakhtin/goph-keeper/gen/go/v1/kv"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/domain/models"
	primary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/ports/primary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/webserver/grpc/interceptors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-faster/errors"
)

var _ kv.KeyValueServiceServer = &KeyValueHandler{}

type KeyValueHandler struct {
	lgr             *zap.Logger
	keyValueUseCase primary_ports.UseCase
	validator       *protovalidate.Validator

	kv.UnimplementedKeyValueServiceServer
}

func NewKeyValueHandler(lgr *zap.Logger, keyValueUseCase primary_ports.UseCase) (*KeyValueHandler, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &KeyValueHandler{
		lgr:             lgr,
		keyValueUseCase: keyValueUseCase,
		validator:       validator,
	}, nil
}

func toModel(ctx context.Context, pb *kv.CrateRequest) (models.KeyValue, error) {
	return models.KeyValue{
		UserID:   ctx.Value(interceptors.ContextKeyUserID).(int),
		Key:      pb.Key,
		Value:    pb.Value,
		Metadata: pb.Metadata,
	}, nil
}

func (h *KeyValueHandler) Create(ctx context.Context, request *kv.CrateRequest) (*emptypb.Empty, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	model, err := toModel(ctx, request)
	if err != nil {
		return nil, err
	}

	err = h.keyValueUseCase.Create(ctx, model)
	if err != nil {
		return nil, errors.Wrap(err, "create key value")
	}

	return &emptypb.Empty{}, nil
}
