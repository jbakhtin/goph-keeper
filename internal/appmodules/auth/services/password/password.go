package password

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
)

type Config interface {
	GetAppKey() string
}

type Service struct {
	cfg Config
	lgr ports.Logger
}

func New(cfg Config, lgr ports.Logger) (*Service, error) {
	return &Service{
		cfg: cfg,
		lgr: lgr,
	}, nil
}

// HashPassword преобразовать незащищенный пароль пользователя в хэш
func (p *Service) HashPassword(password string) (string, error) {
	return p.hashPassword(password)
}

func (p *Service) CheckPassword(password, need string) (bool, error) {
	hashedPassword, err := p.hashPassword(password)
	if err != nil {
		return false, err
	}

	return hashedPassword == need, nil
}

func (p *Service) hashPassword(password string) (string, error) {
	h := hmac.New(sha256.New, []byte(p.cfg.GetAppKey()))
	h.Write([]byte(password))
	dst := h.Sum(nil)

	return fmt.Sprintf("%x", dst), nil
}
