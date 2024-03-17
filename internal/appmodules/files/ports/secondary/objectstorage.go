package secondary

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/files/domain/models"
)

type ObjectStorageRepository interface {
	Save(context.Context, *models.File) error
}
