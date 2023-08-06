package jwt_handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/agaddis02/combocurve-api-go/auth/service_account"
)

const (
	JWTAlgorithm       = "RS256"
	MaxExpirationDelta = 60 * 60 * 24 * 7 // 7 days in seconds
	MinExpirationDelta = 60 * 60          // 1 hour in seconds
)

type JWTHandler struct {
	ServiceAccount      *service_account.ServiceAccount
	Audience            string
	ExpirationTime      time.Duration
	SecondsBeforeExpire time.Duration
}

func NewJWTHandler(sa *service_account.ServiceAccount, secondsBeforeExpire, tokenDuration int, audience string) (*JWTHandler, error) {
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

func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key not in correct format")
	}

	return rsaPrivateKey, nil
}

func ParseRSAPublicKeyFromPrivateKeyPEM(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key not in correct format")
	}

	// Return the public key
	return &rsaPrivateKey.PublicKey, nil
}

func CreateClaims(j *JWTHandler) jwt.Claims {
	now := time.Now()
	return jwt.RegisteredClaims{
		Issuer:    j.ServiceAccount.ClientEmail,
		Subject:   j.ServiceAccount.ClientEmail,
		Audience:  jwt.ClaimStrings{j.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.ExpirationTime)),
	}
}

func (j *JWTHandler) GenerateToken() (string, error) {
	privKey, err := ParseRSAPrivateKeyFromPEM([]byte(j.ServiceAccount.PrivateKey))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(JWTAlgorithm), CreateClaims(j))
	token.Header["kid"] = j.ServiceAccount.PrivateKeyID
	return token.SignedString(privKey)
}

func (j *JWTHandler) IsTokenExpired(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return ParseRSAPublicKeyFromPrivateKeyPEM([]byte(j.ServiceAccount.PrivateKey))
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("invalid token claims")
	}

	fmt.Println(claims.GetIssuer())

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return false, errors.New("Issue getting expiration time: " + err.Error())
	}

	expTime := time.Unix(exp.Unix(), 0)
	return expTime.Before(time.Now().Add(j.SecondsBeforeExpire)), nil
}
