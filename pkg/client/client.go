package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type VGSClient struct {
	token JWT
}

// JWT is a JWT
type JWT struct {
	AccessToken      string `json:"access_token,omitempty"`
	IDToken          string `json:"id_token,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	NotBeforePolicy  int    `json:"not-before-policy,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

func New(ctx context.Context, clientId string, clientSecret string) (*VGSClient, error) {
	var (
		body   = strings.NewReader(`grant_type=client_credentials`)
		client = &http.Client{}
		jwt    = &JWT{}
	)
	req, err := http.NewRequest("POST", "https://auth.verygoodsecurity.com/auth/realms/vgs/protocol/openid-connect/token", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientId, clientSecret)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(jwt)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("token is not valid")
	}

	vc := VGSClient{
		token: JWT{
			AccessToken:      jwt.AccessToken,
			ExpiresIn:        jwt.ExpiresIn,
			RefreshExpiresIn: jwt.RefreshExpiresIn,
			Scope:            jwt.Scope,
			TokenType:        jwt.TokenType,
			NotBeforePolicy:  jwt.NotBeforePolicy,
		},
	}

	return &vc, nil
}

func (v *VGSClient) GetToken() string {
	return v.token.AccessToken
}
