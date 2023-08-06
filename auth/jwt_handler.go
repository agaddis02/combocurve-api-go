package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	JWTAlgorithm       = "RS256"
	MaxExpirationDelta = 60 * 60 * 24 * 7 // 7 days in seconds
	MinExpirationDelta = 60 * 60          // 1 hour in seconds
)

type JWTHandler struct {
	ServiceAccount      *ServiceAccount
	Audience            string
	ExpirationTime      time.Duration
	SecondsBeforeExpire time.Duration
}

func NewJWTHandler(sa *ServiceAccount, secondsBeforeExpire, tokenDuration int, audience string) (*JWTHandler, error) {
	if tokenDuration < MinExpirationDelta || tokenDuration > MaxExpirationDelta {
		return nil, errors.New("tokenDuration should be greater than or equal to MinExpirationDelta and less than or equal to MaxExpirationDelta")
	}
	return &JWTHandler{
		ServiceAccount:      sa,
		Audience:            audience,
		ExpirationTime:      time.Duration(tokenDuration) * time.Second,
		SecondsBeforeExpire: time.Duration(secondsBeforeExpire) * time.Second,
	}, nil
}

func (j *JWTHandler) GenerateToken() (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.GetSigningMethod(JWTAlgorithm), &jwt.StandardClaims{
		Issuer:    j.ServiceAccount.ClientEmail,
		Subject:   j.ServiceAccount.ClientEmail,
		Audience:  j.Audience,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(j.ExpirationTime).Unix(),
	})
	token.Header["kid"] = j.ServiceAccount.PrivateKeyID
	return token.SignedString([]byte(j.ServiceAccount.PrivateKey))
}

func (j *JWTHandler) IsTokenExpired(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.ServiceAccount.PrivateKey), nil
	})
	if err != nil {
		return false, err
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return false, errors.New("invalid token claims")
	}
	expTime := time.Unix(claims.ExpiresAt, 0)
	return expTime.Before(time.Now().Add(j.SecondsBeforeExpire)), nil
}
