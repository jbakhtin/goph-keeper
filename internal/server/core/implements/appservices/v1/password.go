package appservices

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/appservices/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/config"
)

var _ appservices.PasswordAppServiceInterface = &PasswordAppService{}

type PasswordAppService struct {
	cfg config.Config
}

func NewPasswordAppService(cfg config.Config) (*PasswordAppService, error) {
	return &PasswordAppService{
		cfg: cfg,
	}, nil
}

// HashPassword преобразовать незащищенный пароль пользователя в хэш
func (p *PasswordAppService) HashPassword(password string) (string, error) {
	return p.hashPassword(password)
}

func (p *PasswordAppService) CheckPassword(password, need string) (bool, error) {
	hashedPassword, err := p.hashPassword(password)
	if err != nil {
		return false, err
	}

	return hashedPassword == need, nil
}

func (p *PasswordAppService) hashPassword(password string) (string, error) {
	h := hmac.New(sha256.New, []byte(p.cfg.GetAppKey()))
	h.Write([]byte(fmt.Sprintf("%s", password)))
	dst := h.Sum(nil)

	return fmt.Sprintf("%x", dst), nil
}
