package credentials

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-faster/errors"
	"google.golang.org/grpc/credentials"
)

// GetJWTCredentials Функция для создания авторизационных данных на уровне вызовов RPC
func NewJWTCredentials() credentials.PerRPCCredentials {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)
	if err != nil {
		log.Fatal(errors.Wrap(err, "read file to buffer"))
	}

	data, err := reader.ReadBytes('\n')

	var tokens map[string]string
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal file to structure"))
	}

	return jwtCredentials{
		Token: tokens["access_token"],
	}
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
