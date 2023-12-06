package types

import "github.com/golang-jwt/jwt/v5"

type AccessToken string

type TokensPair struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

type UserData struct {
	UserID    ID `json:"user_id,omitempty"`
	SessionID ID `json:"session_id,omitempty"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Data UserData `json:"data,omitempty"`
}

//token, err := jwt.ParseWithClaims(
//	signedToken,
//	&JWTClaim{},
//	func(token *jwt.Token) (interface{}, error) {
//		return []byte(jwtKey), nil
//	},
//)
