package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type VGSClient struct {
	httpClient *uhttp.BaseHttpClient
	token      *JWT
	endpoint   string
}

const (
	applicationJSONHeader     = "application/json"
	applicationFormUrlencoded = "application/x-www-form-urlencoded"
)

func WithBody(body string) uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		var buffer bytes.Buffer
		_, err := buffer.WriteString(body)
		if err != nil {
			return nil, nil, err
		}

		_, headers, err := WithContentTypeFormHeader()()
		if err != nil {
			return nil, nil, err
		}

		return &buffer, headers, nil
	}
}

func WithContentTypeFormHeader() uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		return nil, map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}, nil
	}
}

func New(ctx context.Context, clientId string, clientSecret string) (*VGSClient, error) {
	var jwt = &JWT{}
	uri, err := url.ParseRequestURI("https://auth.verygoodsecurity.com/auth/realms/vgs/protocol/openid-connect/token")
	if err != nil {
		return nil, err
	}

	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli := uhttp.NewBaseHttpClient(httpClient)
	req, err := cli.NewRequest(
		ctx,
		http.MethodPost,
		uri,
		uhttp.WithAcceptJSONHeader(),
		WithBody(`grant_type=client_credentials`),
	)
	req.SetBasicAuth(clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req)
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

func (v *VGSClient) GetOrganizations(ctx context.Context) ([]Organization, error) {
	uri, _ := url.JoinPath(v.endpoint, "/organizations")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", applicationJSONHeader)
	req.Header.Add("Content-Type", applicationFormUrlencoded)
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
	url, _ := url.JoinPath(v.endpoint, "/", orgId, "/users")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", applicationJSONHeader)
	req.Header.Add("Content-Type", applicationFormUrlencoded)
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
