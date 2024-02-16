package credentials

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/go-faster/errors"
	"google.golang.org/grpc/credentials"
)

// GetJWTCredentials Функция для создания авторизационных данных на уровне вызовов RPC
func NewJWTCredentials() (credentials.PerRPCCredentials, error) {
	file, err := os.Open("./config.json") // ToDo: move to config
	if err != nil {
		return nil, errors.Wrap(err, "open secrets")
	}

	reader := bufio.NewReader(file)

	data, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "read bytes to /n")
	}

	var tokens map[string]string
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal secrets to structure")
	}

	return jwtCredentials{
		Token: tokens["access_token"],
	}, nil
}

// jwtCredentials Реализация интерфейса credentials.PerRPCCredentials для передачи токена авторизации в каждом вызове RPC
type jwtCredentials struct {
	Token string
}

func (j jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + j.Token,
	}, nil
}

func (j jwtCredentials) RequireTransportSecurity() bool {
	return false
}
