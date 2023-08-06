package token_refresher

import (
	"github.com/Arch-Energy-Partners/combocurve-api-go/auth/jwt_handler"
	"github.com/Arch-Energy-Partners/combocurve-api-go/auth/service_account"
)

type TokenRefresher struct {
	JwtHandler *jwt_handler.JWTHandler
	Token      string
}

func NewTokenRefresher(sa *service_account.ServiceAccount, secondsBeforeExpire, tokenDuration int, audience string) (*TokenRefresher, error) {
	handler, err := jwt_handler.NewJWTHandler(sa, secondsBeforeExpire, tokenDuration, audience)
	if err != nil {
		return nil, err
	}
	return &TokenRefresher{
		JwtHandler: handler,
	}, nil
}

func (t *TokenRefresher) GetAccessToken() (string, error) {
	if t.Token == "" {
		var err error
		t.Token, err = t.JwtHandler.GenerateToken()
		if err != nil {
			return "", err
		}
	} else {
		isExpired, err := t.JwtHandler.IsTokenExpired(t.Token)
		if err != nil {
			return "", err
		}
		if isExpired {
			t.Token, err = t.JwtHandler.GenerateToken()
			if err != nil {
				return "", err
			}
		}
	}
	return t.Token, nil
}
