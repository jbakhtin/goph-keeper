package types

import "github.com/golang-jwt/jwt/v5"

type TokensPair struct {
	AccessToken  string
	RefreshToken string
}

type UserData struct {
	UserID    int `json:"user_id,omitempty"`
	SessionID int `json:"session_id,omitempty"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Data UserData `json:"data,omitempty"`
}
