package credentials

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-faster/errors"
	"google.golang.org/grpc/credentials"
)

type TokensPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Read() (*TokensPair, error) {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	var tokens TokensPair
	err = json.Unmarshal(content, &tokens)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal file to structure")
	}

	return &tokens, nil
}

func NewRefreshTokenCredentials() credentials.PerRPCCredentials {
	tokens, _ := Read()
	return refreshToken{tokens.RefreshToken}
}

type refreshToken struct {
	Token string
}

func (r refreshToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"refresh-token": r.Token,
	}, nil
}
func (r refreshToken) RequireTransportSecurity() bool {
	return false
}
