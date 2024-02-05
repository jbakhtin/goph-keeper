package password

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
)

type Config interface {
	GetAppKey() string
}

type PasswordAppService struct {
	cfg Config
	lgr ports.Logger
}

func NewPasswordAppService(cfg Config, lgr ports.Logger) (*PasswordAppService, error) {
	return &PasswordAppService{
		cfg: cfg,
		lgr: lgr,
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
	h.Write([]byte(password))
	dst := h.Sum(nil)

	return fmt.Sprintf("%x", dst), nil
}
