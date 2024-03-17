package services

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/files/domain/models"
	primary_ports "github.com/jbakhtin/goph-keeper/internal/appmodules/files/ports/primary"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/appmodules/files/ports/secondary"
)

var _ primary_ports.UseCase = &UseCase{}

type Config struct{}

type UseCase struct {
	cfg                     Config
	lgr                     secondary_ports.Logger
	databaseRepository      secondary_ports.DataBaseRepository
	objectStorageRepository secondary_ports.ObjectStorageRepository
}

func (u *UseCase) Save(ctx context.Context, file *models.File) (*models.File, error) {
	err := u.objectStorageRepository.Save(ctx, file)
	if err != nil {
		return nil, err
	}

	file, err = u.databaseRepository.Create(ctx, file)
	if err != nil {
		return nil, err
	}

	return file, nil
}
