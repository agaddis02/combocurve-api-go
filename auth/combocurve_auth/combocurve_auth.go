package combocurve_auth

import (
	"fmt"

	"github.com/Arch-Energy-Partners/combocurve-api-go/auth/service_account"
	"github.com/Arch-Energy-Partners/combocurve-api-go/auth/token_refresher"
)

const DefaultAudience = "https://api.combocurve.com"

type ComboCurveAuth struct {
	TokenRefresher *token_refresher.TokenRefresher
	ApiKey         string
}

func NewComboCurveAuth(sa *service_account.ServiceAccount, apiKey string, secondsBeforeTokenExpire, tokenDuration int) (*ComboCurveAuth, error) {
	refresher, err := token_refresher.NewTokenRefresher(sa, secondsBeforeTokenExpire, tokenDuration, DefaultAudience)
	if err != nil {
		return nil, err
	}
	return &ComboCurveAuth{
		TokenRefresher: refresher,
		ApiKey:         apiKey,
	}, nil
}

func (c *ComboCurveAuth) GetAuthHeaders() (map[string]string, error) {
	token, err := c.TokenRefresher.GetAccessToken()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"x-api-key":     c.ApiKey,
	}, nil
}
