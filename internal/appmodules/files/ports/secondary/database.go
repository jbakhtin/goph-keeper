package secondary

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/files/domain/models"
)

type DataBaseRepository interface {
	Create(context.Context, *models.File) (*models.File, error)
}
