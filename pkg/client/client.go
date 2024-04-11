package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type VGSClient struct {
	httpClient *http.Client
	token      *JWT
	endpoint   string
}

const (
	applicationJSONHeader     = "application/json"
	applicationFormUrlencoded = "application/x-www-form-urlencoded"
)

func New(ctx context.Context, clientId string, clientSecret string) (*VGSClient, error) {
	var (
		body = strings.NewReader(`grant_type=client_credentials`)
		cli  = &http.Client{}
		jwt  = &JWT{}
	)

	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	url := "https://auth.verygoodsecurity.com/auth/realms/vgs/protocol/openid-connect/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", applicationJSONHeader)
	req.Header.Add("Content-Type", applicationFormUrlencoded)
	req.SetBasicAuth(clientId, clientSecret)

	resp, err := httpClient.Do(req)
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
		httpClient: cli,
		token: &JWT{
			AccessToken:      jwt.AccessToken,
			ExpiresIn:        jwt.ExpiresIn,
			RefreshExpiresIn: jwt.RefreshExpiresIn,
			Scope:            jwt.Scope,
			TokenType:        jwt.TokenType,
			NotBeforePolicy:  jwt.NotBeforePolicy,
		},
		endpoint: "https://accounts.apps.verygoodsecurity.com",
	}

	return &vc, nil
}

func (v *VGSClient) GetToken() string {
	return v.token.AccessToken
}

func (v *VGSClient) InitHttpClient() {
	v.httpClient = &http.Client{}
}

func (v *VGSClient) GetOrganizations(ctx context.Context) ([]Organization, error) {
	v.InitHttpClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.endpoint+"/organizations", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+v.GetToken())
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var organizationsAPIData organizationsAPIData
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &organizationsAPIData)
	if err != nil {
		return nil, err
	}

	var organizations []Organization
	for _, org := range organizationsAPIData.Data {
		organizations = append(organizations, Organization{
			Id:        org.Id,
			Name:      org.Attributes.Name,
			State:     org.Attributes.State,
			CreatedAt: org.Attributes.CreatedAt,
			UpdatedAt: org.Attributes.UpdatedAt,
		})
	}

	return organizations, nil
}

func (v *VGSClient) GetOrganizationUsers(ctx context.Context, orgId string) ([]OrganizationUser, error) {
	v.InitHttpClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.endpoint+"/"+orgId+"/users", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+v.GetToken())
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var organizationUsersAPIData organizationUsersAPIData
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &organizationUsersAPIData)
	if err != nil {
		return nil, err
	}

	var users []OrganizationUser
	for _, userAPI := range organizationUsersAPIData.Data {
		users = append(users, OrganizationUser{
			Id:        userAPI.Id,
			Name:      userAPI.Attributes.Name,
			Email:     userAPI.Attributes.EmailAddress,
			CreatedAt: userAPI.Attributes.CreatedAt,
			UpdatedAt: userAPI.Attributes.UpdatedAt,
		})
	}

	return users, nil
}
