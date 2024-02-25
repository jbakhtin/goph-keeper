package primary

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/files/domain/models"
)

type UseCase interface {
	Save(context.Context, *models.File) (*models.File, error)
}
